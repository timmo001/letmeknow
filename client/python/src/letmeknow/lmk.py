"""Client for communication with LetMeKnow."""

from __future__ import annotations

import asyncio
from dataclasses import dataclass
from importlib import metadata
from socket import gaierror
from typing import TYPE_CHECKING, Any, Final

from aiohttp import (
    ClientConnectionError,
    ClientSession,
    ClientWebSocketResponse,
    WSServerHandshakeError,
)
from yarl import URL

from .exceptions import LMKConnectionError, LMKNotConnectedError
from .models import (
    LMKClientType,
    LMKNotification,
    LMKWSNotification,
    LMKWSRegister,
    LMKWSRequestType,
    LMKWSResponseError,
    LMKWSResponseSuccess,
)

if TYPE_CHECKING:
    from collections.abc import Callable

    from typing_extensions import Self

VERSION: Final[str] = metadata.version(__package__)


@dataclass
class LMKClient:
    """Client to communicate with LetMeKnow."""

    lmk_host: str
    lmk_port: int
    lmk_client_type: LMKClientType
    lmk_user_id: str

    session: ClientSession | None = None
    request_timeout: int = 10
    _close_session: bool = False
    _ws: ClientWebSocketResponse | None = None

    async def __aenter__(self) -> Self:
        """Async enter.

        Returns
        -------
            The LMKClient object.

        """
        return self

    async def __aexit__(self, *_exc_info: object) -> None:
        """Async exit.

        Args:
        ----
            _exc_info: Exec type.

        """
        await self.close()

    async def _ws_send(
        self,
        data: dict[str, Any],
    ) -> LMKWSResponseSuccess | LMKWSResponseError:
        """Send data to the websocket server.

        Args:
        ----
            data: Data to send.

        """
        if self._ws is None:
            raise LMKNotConnectedError

        await self._ws.send_json(data)

        async with asyncio.timeout(self.request_timeout):
            response = await self._ws.receive()

        # Get the response data
        response_data: dict[str, Any] = response.json()

        if response.type == 1:
            return LMKWSResponseSuccess.from_dict(response_data)
        return LMKWSResponseError.from_dict(response_data)

    async def close(self) -> None:
        """Close open client session."""
        if self.session and self._close_session:
            await self.session.close()

    async def ws_connect(self) -> bool:
        """Connect to LetMeKnow websocket server.

        Returns
        -------
            True if the connection was successful.

        """
        url = URL.build(
            scheme="ws",
            host=self.lmk_host,
            port=self.lmk_port,
        ).joinpath("/websocket")

        headers = {
            "User-Agent": f"LMKClientPy/{VERSION}",
        }

        if self.session is None:
            self.session = ClientSession()
            self._close_session = True

        try:
            async with asyncio.timeout(self.request_timeout):
                self._ws = await self.session.ws_connect(
                    url=url,
                    headers=headers,
                    heartbeat=30,
                )
        except (
            asyncio.TimeoutError,
            WSServerHandshakeError,
            ClientConnectionError,
            gaierror,
        ) as error:
            raise LMKConnectionError from error

        return True

    async def ws_register(self) -> LMKWSResponseSuccess | LMKWSResponseError:
        """Register with the websocket server.

        Returns
        -------
            The response from the websocket server.

        """
        return await self._ws_send(
            LMKWSRegister(
                type=LMKWSRequestType.REGISTER,
                user_id=self.lmk_user_id,
            ).to_dict()
        )

    async def ws_send_notification(
        self,
        notification: LMKNotification,
    ) -> LMKWSResponseSuccess | LMKWSResponseError:
        """Send a notification to the websocket server.

        Args:
        ----
            notification: Notification to send.

        Returns:
        -------
            The response from the websocket server.

        """
        return await self._ws_send(
            LMKWSNotification(
                type=LMKWSRequestType.NOTIFICATION,
                data=notification,
                targets=[],
            ).to_dict()
        )

    async def ws_listen_for_notifications(
        self,
        cb: Callable[[LMKNotification], None],
    ) -> None:
        """Listen for notifications.

        Args:
        ----
            cb: Callback to call when a notification is received.

        """
        if self._ws is None:
            raise LMKNotConnectedError

        async for message in self._ws:
            if message.type == 1:
                response_data: dict[str, Any] = message.json()
                notification = LMKWSNotification.from_dict(response_data)
                cb(notification.data)
            else:
                break
