use serde::{Deserialize, Serialize};

use crate::routine::Routine;
use crate::CONFIG_FILE;
use std::io::prelude::*;

use crate::Cache;
use crate::{read_config, write_config};

use chrono;
use chrono_humanize::HumanTime;
use std::fs;
use std::io;
use std::time::{Duration, SystemTime};
#[derive(Debug, Serialize, Deserialize)]
pub struct TaskControl {
    routines: Vec<Routine>,
    cache: Cache,
}

impl TaskControl {
    /// execute all routines
    pub fn execute_all(&mut self) -> io::Result<()> {
        debug!("run TaskControl.execute_all");
        for routine in self.routines.iter() {
            let key = &routine.name;
            let mut should_execute = true;
            let now = SystemTime::now();
            let mut updated_secs_before: Option<i64> = None;
            if let Some(entry) = self.cache.entries.get_mut(key) {
                let one_minute = Duration::new(60, 0);
                let min_delay = one_minute * routine.interval_minute;
                let last_update = entry.last_mod;

                let earliest_time_to_update = last_update + min_delay;
                should_execute = earliest_time_to_update < now;
                if !should_execute {
                    if let Ok(dura) = now.duration_since(last_update) {
                        let secs = dura.as_secs() as i64;
                        updated_secs_before = Some(secs);
                    }
                }
            }
            if should_execute {
                debug!("Should update");
                if let Ok(_) = routine.execute() {
                    self.cache.update(key)?;
                    debug!("Update cache and flush out..");
                    self.cache.export()?;
                }
            } else {
                if let Some(secs_before) = updated_secs_before {
                    let dura = chrono::Local::now() - chrono::Duration::seconds(secs_before);
                    // @TODO: show currently running command
                    info!("{}", format!("Have update {}", HumanTime::from(dura)));
                }
                debug!("Cache entry exist. jump to the next");
            }
        }
        self.cache.export()?;
        Ok(())
    }

    pub fn replace_cache(&mut self, cache: Cache) {
        self.cache = cache;
    }

    fn new() -> TaskControl {
        let routines: Vec<Routine> = Vec::new();
        return TaskControl {
            routines,
            cache: Cache::new(),
        };
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
        TaskControl {
            routines,
            cache: Cache::new(),
        }
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
