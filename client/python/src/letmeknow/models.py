"""Python client for LetMeKnow."""

from __future__ import annotations

from dataclasses import Field, dataclass
from enum import StrEnum
from typing import Any, Self
from uuid import uuid4


class LMKClientType(StrEnum):
    """Enum of client type."""

    CLIENT = "client"
    HEADLESS = "headless"


@dataclass(slots=True)
class LMKNotification:
    """Notification."""

    type: str | None = None
    title: str | None = None
    subtitle: str | None = None
    content: str | None = None
    image: LMKNotificationImage | None = None

    @classmethod
    def from_dict(cls, result: dict[str, Any]) -> Self:
        """Initialize from a dict."""
        return cls(
            type=result["type"],
            title=result["title"],
            subtitle=result["subtitle"],
            content=result["content"],
            image=LMKNotificationImage.from_dict(result["image"]),
        )

    @classmethod
    def to_dict(cls) -> dict[str, Any]:
        """Convert to a dict."""
        return {
            "type": cls.type,
            "title": cls.title,
            "subtitle": cls.subtitle,
            "content": cls.content,
            "image": cls.image.to_dict() if cls.image else None,
        }


@dataclass(slots=True)
class LMKNotificationImage:
    """Notification image."""

    url: str
    height: float | None = None
    width: float | None = None

    @classmethod
    def from_dict(cls, result: dict[str, Any]) -> Self:
        """Initialize from a dict."""
        return cls(
            url=result["url"],
            height=result["height"],
            width=result["width"],
        )

    @classmethod
    def to_dict(cls) -> dict[str, Any]:
        """Convert to a dict."""
        return {
            "url": cls.url,
            "height": cls.height,
            "width": cls.width,
        }


class LMKWSRequestType(StrEnum):
    """Enum of websocket request type."""

    REGISTER = "register"
    NOTIFICATION = "notification"


class LMKWSResponseType(StrEnum):
    """Enum of websocket response type."""

    ERROR = "error"
    NOTIFICATION_SENT = "notificationSent"
    REGISTER = "register"


@dataclass(slots=True)
class LMKWSRegister:
    """Websocket register client."""

    type: LMKWSRequestType
    user_id: str

    def generate_user_id(self) -> str:
        """Generate a user ID."""
        return f"{self.type}-{uuid4()!s}"

    @classmethod
    def from_dict(cls, result: dict[str, Any]) -> Self:
        """Initialize from a dict."""
        return cls(
            type=LMKWSRequestType(result["type"]),
            user_id=result["userID"],
        )

    @classmethod
    def to_dict(cls) -> dict[str, Any]:
        """Convert to a dict."""
        return {
            "type": cls.type,
            "userID": cls.user_id,
        }


@dataclass(slots=True)
class LMKWSNotification:
    """Websocket notification."""

    type: LMKWSRequestType
    data: LMKNotification
    targets: list[str]

    @classmethod
    def from_dict(cls, result: dict[str, Any]) -> Self:
        """Initialize from a dict."""
        return cls(
            type=LMKWSRequestType(result["type"]),
            data=LMKNotification.from_dict(result["data"]),
            targets=result["targets"],
        )

    @classmethod
    def to_dict(cls) -> dict[str, Any]:
        """Convert to a dict."""
        return {
            "type": cls.type,
            "data": cls.data.to_dict(),
            "targets": cls.targets,
        }


@dataclass(slots=True)
class LMKWSResponseError:
    """Websocket response error."""

    type: LMKWSResponseType
    message: str
    error: str

    @classmethod
    def from_dict(cls, result: dict[str, Any]) -> Self:
        """Initialize from a dict."""
        return cls(
            type=LMKWSResponseType(result["type"]),
            message=result["message"],
            error=result["error"],
        )


@dataclass(slots=True)
class LMKWSResponseSuccess:
    """Websocket response success."""

    type: LMKWSResponseType
    succeeded: bool
    message: str

    @classmethod
    def from_dict(cls, result: dict[str, Any]) -> Self:
        """Initialize from a dict."""
        return cls(
            type=LMKWSResponseType(result["type"]),
            succeeded=result["succeeded"],
            message=result["message"],
        )
