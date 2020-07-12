# daily-runner

`daily-runner` is an application that allows you to run other command line
applications in a "daily" manner even when a computer is not running 24/7. The
command will be run if more than 24 hours have passed since the last execution.
A preferred run time can also be set that overrides the 24h mode and will
always run if the system is up at this preferred time. It is designed to easily
run as a non system application with minimal installation by a typical non-root
desktop linux user.

## Quick start

* Grab the latest binary release from
  [releases](https://github.com/gaganis/daily-runner/releases) and put it
somewhere in your home directory. `$HOME/bin` is a traditional location for
user-level executables. Alternative you can get the source and use `go install
daily-runner` this will put the executatable in your `$GOPATH/bin` folder.

* Optionally run `daily-runner` with your command line flags from the command
  line for example for my case

```
gaganis@i7:~$ $GOPATH/bin/daily-runner -profile backup_i7_s3 -command "duply backup_i7_s3 backup" -preferredTime "03:00:00" 
Starting daily-runner with configuration:
{Profile:backup_i7_s3 Command:duply backup_i7_s3 backup Interval:4m0s HasPreferredRunTime:true PreferredRunTime:03:00:00}
Please see logs at:  /home/gaganis/.local/share/daily-runner/backup_i7_s3/log/daily-runner.log
```

* To make `daily-runner` start when you login in your system you can create a
  file in `$HOME/.config/autostart` named `<something>.desktop` For example my
desktop file for the above profile & command `daily-runner-backup-s3.desktop`:

``` 
[Desktop Entry]
Type=Application
Exec=/home/gaganis/GolandProjects/bin/daily-runner -profile backup_i7_s3 -command "duply backup_i7_s3 backup" -preferredTime "03:00:00" 
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name[en_US]=daily runner duply s3 backup
Name=daily runner duply s3 backup
Comment[en_US]=
Comment=
```

* Alternatively you can you use the gnome UI app `Startup Applications`
  accessible from you distro's launcher

## Usage

```
Usage of /home/gaganis/GolandProjects/bin/daily-runner:
  -command string
        The command that runner will execute (default "echo 'daily-runner has run echo printing this text'")
  -interval duration
        The interval that daily-runner will use to check if it needs to run. Can accept values acceptable to golang time.ParseDuration function (default 4m0s)
  -preferredTime string
        Set a preferred time for the runner to run command. This time overrides the daily logic and the command will always run if the system is up at that time.
  -profile string
        Profile to use. Defaults to 'default' (default "default")
```

## Features

#### 24h run algorithm

Every time `daily-runner` executes the command it tracks the time and if more
than 24 hours have passed since the last run.

#### log files

Logs are saved in `$HOME/.local/share/daily-runner/<profile>/log` in the
following files:
 * `daily-runner.log` 
 * `command-output.log` 

#### Preferred run time

A preferred run time can be set via the command line flag `-preferredTime` ie
01:00:00, if daily-runner is running at this time the command will be executed
irrespective of when the last run was. 

Please note that daily-run does not attempt to run the backup at exactly the
time. It can be up to time adding the interval setting. For the default 4min
interval this means that the command can run from 01:00:00 up to 01:04:10.

#### Configurable check interval

`daily-runner` checks in a loop whether the command needs to run. This loop by
default sleeps 4min between checks to conserve power. This can be adjusted via
command line flag `-interval`.

#### Profiles A profile can be set via the `-profile` command line flag. This
profile allows multiple instances to run for different commands independently
each with it's it's own logging and data about runs. If no profile is set
`default` is used.

### Single process running via pid file locking

Locking is implemented via a pid file so that no more than one instance of
`daily-runner` is running for each profile.

## Motivation

I have been using [duply](https://duply.net/index.php/Main_Page) at work with
great success, even used it once to migrate to a new computer, and I wanted to
also setup backups for my home computers. 

At work was running it via cron at lunchtime but for my home computers which I
use randomly & sporadically I needed to find a different way to run it. I
looked at fcron & anacron and I found them hard to install requiring system
changes(fcron an compilation from sources) plus a steep learning & setup curve
to reach my use case. 

Gnome's [deja dup](https://wiki.gnome.org/Apps/DejaDup) is also an app that
tries to solve the overall problem of home/personal computer backup and has the
feature of "daily", "weekly" running of backup for systems that are not always
online.  It's design imposes many defaults, so I was not able to set it up
successfully for my use cases. 
