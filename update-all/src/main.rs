use clap::{self, App, Arg};
// use serde::{Deserialize, Serialize};
// use serde_yaml;
// use std::format;
use std::fs::{self, OpenOptions};
use std::io;
use std::io::prelude::*;
use std::path::Path;

#[macro_use]
extern crate log;

use log::Level;
use loggerv;

mod cache;
mod routine;
mod task_control;
use cache::Cache;
use task_control::TaskControl;

/// A file that store commands to execute
static CONFIG_FILE: &str = ".update_all.config.yaml";
static CACHE_FILE: &str = ".update_time.json";

fn config_exists() -> bool {
    Path::new(CONFIG_FILE).exists()
}

fn read_config() -> io::Result<String> {
    // @TODO: create config if not exists
    let raw_config: String = fs::read_to_string(CONFIG_FILE)?;
    Ok(raw_config)
}

/// Write config to "existing" config file
fn write_config(str: String) -> io::Result<()> {
    let mut file = fs::OpenOptions::new()
        .append(true)
        .open(CONFIG_FILE)
        .unwrap();
    file.write(str.as_bytes()).expect("Cannot write");
    Ok(())
}

#[derive(Debug)]
struct CliConfig {
    force_all: bool,
    debug: bool,
    create: bool,
}

impl CliConfig {
    fn new(matches: clap::ArgMatches) -> CliConfig {
        let force_all = matches.is_present("force-all");
        let debug = matches.is_present("debug");
        let create = matches.is_present("create");
        return CliConfig {
            force_all,
            debug,
            create,
        };
    }
}

fn main() -> Result<(), io::Error> {
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
        .arg(
            Arg::with_name("create")
                .long("create")
                .help("Create Default config file if not exists, would return after file has been created.")
                .takes_value(false),
        );

    let matches = app.get_matches();
    let config = CliConfig::new(matches);

    if config.debug {
        loggerv::init_with_level(Level::Trace).unwrap();
    } else {
        loggerv::init_with_level(Level::Info).unwrap();
    }
    info!("Start to execute routines");

    let cfg_exists = config_exists();
    if !cfg_exists {
        debug!("Config file doesn't exists");
        if config.create {
            {
                // ensure file exists
                OpenOptions::new()
                    .write(true)
                    .create_new(true)
                    .open(CONFIG_FILE)
                    .unwrap();
            }
            let default_tasks = TaskControl::default_template();
            default_tasks
                .export_routine_append()
                .expect("cannot export routine");
            println!("Create a default config file: {}", CONFIG_FILE);
        } else {
            //@TODO: add colors
            println!("");
            println!("Config file not exist");
            println!("Run command with --create");
            println!("To create a config file");
        }
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
    }

    info!("Start to execute routines");
    taskctl.execute_all().expect("Cannot execute command");
    Ok(())
}
