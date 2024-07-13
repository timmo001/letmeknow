#[cfg(target_os = "macos")]
#[macro_use]
extern crate objc;

use tauri::{
    menu::{MenuBuilder, MenuItemBuilder, PredefinedMenuItem},
    tray::{MouseButton, MouseButtonState, TrayIconEvent},
    Manager,
};

#[tauri::command]
async fn set_window(app: tauri::AppHandle) -> Result<(), String> {
    println!("Setting window properties");
    
    // Get the main window
    let window = app.get_webview_window("main").unwrap();

    let _ = window.clone().with_webview(move |_webview| {
        // Allow clickthrough on the window (macOS)
        #[cfg(target_os = "macos")]
        unsafe {
            let () = msg_send![webview.ns_window(), setIgnoresMouseEvents: true];
        }

        // Allow clickthrough on the window (Windows)
        #[cfg(target_os = "windows")]
        unsafe {
            let hwnd = window.hwnd().unwrap().0;
            let hwnd = windows::Win32::Foundation::HWND(hwnd);
            use windows::Win32::UI::WindowsAndMessaging::*;
            let nindex = GWL_EXSTYLE;
            let style = WS_EX_APPWINDOW
                | WS_EX_COMPOSITED
                | WS_EX_LAYERED
                | WS_EX_TRANSPARENT
                | WS_EX_TOPMOST;
            let _pre_val = SetWindowLongA(hwnd, nindex, style.0 as i32);
        }
    });

    Ok(())
}

pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .invoke_handler(tauri::generate_handler![set_window])
        .setup(|app| {
            {
                // Setup tray menu
                let separator = PredefinedMenuItem::separator(app)?;
                let show_settings =
                    MenuItemBuilder::with_id("show_settings", "Open settings").build(app)?;
                let exit = MenuItemBuilder::with_id("exit", "Exit").build(app)?;

                let menu = MenuBuilder::new(app)
                    .items(&[&show_settings, &separator, &exit])
                    .build()?;

                // Setup tray icon
                let tray = app.tray_by_id("main").unwrap();
                tray.set_menu(Some(menu))?;
                tray.on_tray_icon_event(|_tray, event| {
                    if let TrayIconEvent::Click {
                        button: MouseButton::Left,
                        button_state: MouseButtonState::Up,
                        ..
                    } = event
                    {
                        // TODO: Send a message to the webview
                    }
                });
                tray.on_menu_event(move |_app_handle, event| match event.id().as_ref() {
                    "show_settings" => {
                        // TODO: Create a settings window
                    }
                    "exit" => {
                        std::process::exit(0);
                    }
                    _ => {}
                });

                // Get the main window
                let window = app.get_webview_window("main").unwrap();

                // Open devtools on startup
                #[cfg(debug_assertions)] // Only include this code on debug builds
                {
                    window.open_devtools();
                };
            }
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
