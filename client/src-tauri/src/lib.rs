use tauri::{
    menu::{MenuBuilder, MenuItemBuilder, PredefinedMenuItem},
    tray::{MouseButton, MouseButtonState, TrayIconEvent},
    Manager,
};

pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .setup(|app| {
            #[cfg(debug_assertions)] // Only include this code on debug builds
            {
                // Open devtools on start
                let window = app.get_webview_window("main").unwrap();
                window.open_devtools();
                window.close_devtools();

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
                tray.on_tray_icon_event(|tray, event| {
                    if let TrayIconEvent::Click {
                        button: MouseButton::Left,
                        button_state: MouseButtonState::Up,
                        ..
                    } = event
                    {
                        // TODO: Send a message to the webview
                    }
                });
                tray.on_menu_event(|_, event| match event.menu_item_id {
                    "show_settings" => {
                        // TODO: Create a settings window
                    }
                    "exit" => {
                        app.exit(0);
                    }
                });
            }
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
