use log::info;
use serde::{Deserialize, Serialize};

use crate::shared::get_data_path;

#[derive(Debug, Serialize, Deserialize)]
pub struct Settings {
    pub autostart: bool,
    pub log_level: String,
    pub server: SettingsServer,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SettingsServer {
    pub host: String,
    pub port: i64,
}

fn create_settings() -> Settings {
    // Create settings from {config_path}\settings.json
    let settings = Settings {
        autostart: false,
        log_level: "INFO".to_string(),
        server: SettingsServer {
            host: "localhost".to_string(),
            port: 8080,
        },
    };

    // Create settings string
    let settings_string = serde_json::to_string(&settings).unwrap();

    info!("Creating settings file: {}", settings_string);

    // Write settings to {config_path}\settings.json
    let settings_path = format!("{}/settings.json", get_data_path());
    std::fs::write(settings_path, settings_string).unwrap();

    settings
}

#[tauri::command]
pub fn get_settings() -> Settings {
    // Read settings from {config_path}\settings.json
    let settings_path = format!("{}/settings.json", get_data_path());
    if !std::path::Path::new(&settings_path).exists() {
        return create_settings();
    }

    let settings = std::fs::read_to_string(settings_path);
    if settings.is_err() {
        return create_settings();
    }
    let settings = serde_json::from_str(&settings.unwrap());
    if settings.is_err() {
        return create_settings();
    }

    settings.unwrap()
}

#[tauri::command]
pub fn update_settings(settings: Settings) -> Result<(), String> {
    // Write settings to {config_path}\settings.json
    let settings_string = serde_json::to_string(&settings).unwrap();
    let settings_path = format!("{}/settings.json", get_data_path());
    std::fs::write(settings_path, settings_string).unwrap();

    Ok(())
}
