use serde::{Deserialize, Serialize};

use crate::routine::Routine;
use crate::CONFIG_FILE;
use std::io::prelude::*;

use crate::Cache;
use crate::{read_config, write_config};
use chrono::{self, Duration};
use std::fs;
use std::io;

#[derive(Debug, Serialize, Deserialize)]
pub struct TaskControl {
    routines: Vec<Routine>,
}

impl TaskControl {
    /// execute all routines
    pub fn execute_all(&self) -> io::Result<()> {
        debug!("run TaskControl.execute_all");
        for routine in self.routines.iter() {
            // println!("cache exists? :{:?}", Cache::exists(&routine.name));
            // if Cache::exists(&routine.name) {
            //     let last_update = Cache::last_update_utc(&routine.name);
            //     let min_delay = chrono::Duration::minutes(routine.interval_minute);
            //     let next_run = last_update + min_delay;
            //     if next_run <

            //     println!("last update time :{:?}", time);
            // io::stdout().flush()?;
            // } else {
            //     routine.execute()?;
            // }
            // Cache::update(&r.name)?;
        }
        Ok(())
    }

    fn new() -> TaskControl {
        let routines: Vec<Routine> = Vec::new();
        return TaskControl { routines };
    }
    fn add_routine(&mut self, routine: Routine) {
        self.routines.push(routine);
    }

    /// Create a default TaskControl for template
    pub fn default_template() -> TaskControl {
        let mut taskctl = TaskControl::new();
        taskctl.add_routine(Routine::new(
            60,
            "ls".to_string(),
            vec!["-a".into(), "-l".into()],
        ));
        taskctl.add_routine(Routine::new(60, "echo".to_string(), vec!["Good".into()]));
        {
            // prepend file with comments
            let mut file = fs::OpenOptions::new()
                .write(true)
                .create(true)
                .open(CONFIG_FILE)
                .unwrap();
            // let mut file = fs::File::create(CONFIG_FILE).expect("cannot create file");
            file.write("# Change lines below for writinrg files\n".as_bytes())
                .expect("cannot write to file");
        }
        return taskctl;
    }

    /// initialize struct by reading default config file
    pub fn from_cfg_file() -> TaskControl {
        let cfg = read_config().expect("Cannot read config file");

        let routines: Vec<Routine> = serde_yaml::from_str(&cfg).expect("Cannot parse data");
        TaskControl { routines }
    }

    /// Append current routines into existing config file
    pub fn export_routine_append(&self) -> io::Result<()> {
        let cfg = match serde_yaml::to_string(&self.routines) {
            Ok(s) => s,
            Err(e) => return Err(io::Error::new(io::ErrorKind::InvalidData, e)),
        };
        write_config(cfg).expect("cannot write config");
        Ok(())
    }
}
