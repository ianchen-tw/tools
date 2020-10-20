use serde::{Deserialize, Serialize};
use std::io;
use std::process::Command;

#[derive(Debug, Serialize, Deserialize)]
pub struct Routine {
    pub interval_minute: i64,
    pub name: String,
    pub args: Vec<String>,
}

impl Routine {
    pub fn new(interval: i64, name: String, args: Vec<String>) -> Routine {
        return Routine {
            interval_minute: interval,
            name: name,
            args: args,
        };
    }
    pub fn execute(&self) -> Result<(), io::Error> {
        let mut cmd = Command::new(self.name.clone());
        for arg in &self.args {
            cmd.arg(arg.clone());
        }
        log::info!(
            "{}",
            format!("Running Command : {} {:?}", self.name, self.args)
        );
        // @TODO: provide an slient option
        let mut child = cmd
            .spawn()
            .expect(&format!("Failed to execute command: {:?}", self.name));
        let _ecode = child.wait().expect("Failed to wait on child");
        Ok(())
    }
}
