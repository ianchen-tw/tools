use std::collections::hash_map::DefaultHasher;
use std::hash::{Hash, Hasher};
use std::io;
use std::process::Command;

use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Routine {
    pub interval_minute: u32,
    pub name: String,
    pub args: Vec<String>,
}

impl Hash for Routine {
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.name.hash(state);
        self.args.hash(state);
    }
}

impl Routine {
    pub fn new(interval: u32, name: String, args: Vec<String>) -> Routine {
        return Routine {
            interval_minute: interval,
            name: name,
            args: args,
        };
    }
    pub fn clone(&self) -> Routine {
        Routine {
            interval_minute: self.interval_minute,
            name: self.name.clone(),
            args: self.args.clone(),
        }
    }
    pub fn execute(&self, dry_run: bool) -> Result<(), io::Error> {
        let mut cmd = Command::new(self.name.clone());
        for arg in &self.args {
            cmd.arg(arg.clone());
        }
        log::info!(
            "{}",
            format!("Running Command : {} {:?}", self.name, self.args)
        );
        // @TODO: provide an slient option
        if !dry_run {
            let mut child = cmd
                .spawn()
                .expect(&format!("Failed to execute command: {:?}", self.name));
            let _ecode = child.wait().expect("Failed to wait on child");
        }
        Ok(())
    }
    pub fn get_hash(&self) -> u64 {
        let mut s = DefaultHasher::new();
        self.hash(&mut s);
        s.finish()
    }
}
