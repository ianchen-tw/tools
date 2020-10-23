#[macro_use]
extern crate log;

use std::fs::{self, OpenOptions};
use std::io;
use std::io::prelude::*;
use std::process::Command;

use clap::{self, App, Arg};
use colored::*;
use home::home_dir;
use log::Level;
use loggerv;

mod cache;
mod routine;
mod task_control;
use cache::Cache;
use task_control::TaskControl;

// @TODO: remove unsafe use of static mut
// this could be the bug for multithreading
static mut IS_DEBUG: bool = false;
static mut NO_HOME: bool = false;

/// Return the path of different config
fn base_path() -> std::path::PathBuf {
    let mut result_path = std::path::PathBuf::new();
    result_path.push(".");
    unsafe {
        // the base_path in DEBUG mode is default to be current dir.
        if !IS_DEBUG && !NO_HOME {
            if let Some(mut home_path) = home_dir() {
                home_path.push(".config");
                home_path.push("update-all");
                result_path = home_path;
            }
        }
    }
    // make sure result_path exists
    if !result_path.exists() {
        info!("Create base folder: {:#?}", result_path);
        fs::create_dir_all(&result_path).unwrap();
    }
    result_path
}

/// Return the path of the config file.
/// World prefix with home folder, if exists.
fn config_path() -> std::path::PathBuf {
    let mut result = base_path();
    result.push("update-all.config.yaml");
    result
}

/// Return the path of the cache file.
/// Would prefix with home folder, if exists.
fn cache_path() -> std::path::PathBuf {
    let mut result = base_path();
    result.push("update-all.cache.json");
    result
}

fn config_exists() -> bool {
    config_path().exists()
}

fn read_config() -> io::Result<String> {
    let raw_config: String = fs::read_to_string(config_path())?;
    Ok(raw_config)
}

/// Write config to "existing" config file
fn write_config(str: String) -> io::Result<()> {
    let mut file = fs::OpenOptions::new()
        .append(true)
        .open(config_path())
        .unwrap();
    file.write(str.as_bytes()).expect("Cannot write");
    Ok(())
}

#[derive(Debug)]
struct CliConfig {
    force_all: bool,
    debug: bool,
    create: bool,
    dry: bool,
    nohome: bool,
    edit: bool,
}

impl CliConfig {
    fn new(matches: clap::ArgMatches) -> CliConfig {
        let force_all = matches.is_present("force-all");
        let debug = matches.is_present("debug");
        let create = matches.is_present("create");
        let dry = matches.is_present("dry");
        let nohome = matches.is_present("nohome");
        let edit = matches.is_present("edit");
        return CliConfig {
            force_all,
            debug,
            create,
            dry,
            nohome,
            edit,
        };
    }
}

/// Use vscode to open the config file
fn open_editing_config() {
    // @TODO: support default editor
    warn!("open config file");
    let mut cmd = Command::new("code");
    cmd.arg(config_path());

    let mut child = cmd.spawn().expect("Failed to spwan process");
    let _ecode = child.wait().expect("Failed to edit command");
}

fn main() -> Result<(), io::Error> {
    // @TODO: add subcommand for delete cache
    let app = App::new("update-all")
        .version("0.1")
        .about("Run your commands on daily basis")
        .author("Ian Chen")
        .arg(
            Arg::with_name("force-all")
                .short("f")
                .long("force-all")
                .takes_value(false),
        )
        .arg(
            Arg::with_name("debug")
                .short("d")
                .long("debug")
                .takes_value(false),
        )
        // @TODO: turn create argument into a subcommand
        .arg(
            Arg::with_name("create")
                .long("create")
                .help("Create Default config file if not exists, would return after file has been created.")
                .takes_value(false),
        )
        // @TODO: turn edit argument into a subcommand
        .arg(
            Arg::with_name("edit")
            .long("edit")
            .help("open the config file for editing")
            .takes_value(false),
        )
        .arg(
            Arg::with_name("nohome")
            .long("nohome")
            .help("read/write configs in current dir, (not from `~/.config`)")
            .takes_value(false),
        )
        .arg(
            Arg::with_name("dry")
            .long("dry")
            .help("Use dry run (Not executing the command).")
            .takes_value(false),
        );
    let config = CliConfig::new(app.get_matches());

    if config.debug {
        loggerv::init_with_level(Level::Trace).unwrap();
    } else {
        loggerv::init_with_level(Level::Info).unwrap();
    }
    unsafe {
        IS_DEBUG = config.debug;
        NO_HOME = config.nohome;
    }
    info!("update-all Started");

    let cfg_exists = config_exists();
    if !cfg_exists {
        debug!("Cannot find config file");
        if config.create {
            {
                // ensure file exists
                OpenOptions::new()
                    .write(true)
                    .create_new(true)
                    .open(config_path())
                    .unwrap();
            }
            let default_tasks = TaskControl::write_default_template();
            default_tasks
                .export_routine_append()
                .expect("cannot export routine");
            println!("Create a default config file: {:?}", config_path().to_str());
            if config.edit {
                // edit a newly create config file
                open_editing_config();
            } else {
                // info!("{}", "please use".green()  to edit your config file".red());
                info!("{}", "please use:".green());
                info!(
                    "{}{}",
                    " code ".yellow(),
                    config_path().to_str().unwrap().yellow()
                );
                info!("{}", "to edit your config file".green());
            }
        } else {
            warn!("{}", "Config file doesn't exists".red());
            warn!("{}", "Run command with:".red());
            warn!("{}", "   --create");
            warn!("{}", "to create a config file".red());
            warn!(
                "{}{}{}{}{}",
                "(or use ".red(),
                "--edit",
                " with ".red(),
                "--create",
                " to open the config file)".red()
            );
        }
        return Ok(());
    }
    if config.edit {
        // editing existing config file
        open_editing_config();
        return Ok(());
    }
    info!("Load config from file");
    let mut taskctl = TaskControl::from_cfg_file();
    debug!("{}", format!("Import routines : {:#?}", taskctl));

    if config.force_all {
        info!("invalidate cache directory");
        Cache::remove_file().unwrap();
    } else if Cache::could_load_from_file() {
        debug!("Load cache from file");
        let cache = Cache::from_cache_file();
        debug!("{}", format!("Cache: {:#?}", cache));
        taskctl.replace_cache(cache);
    } else {
        // Cache might be corrupted or simply not exists
        Cache::remove_file().unwrap();
    }

    info!("Start to execute routines");
    taskctl
        .execute_all(config.dry)
        .expect("Cannot execute command");
    info!("Success: Routines executed");
    Ok(())
}
