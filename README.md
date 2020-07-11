# daily-runner

`daily-runner` is an application that allows you to run other command line
applications in a "daily" manner even when a computer is not running 24/7. The
command will be run if more than 24 hours have passed since the last execution.
A preferred run time can also be set that overrides the 24h mode and will always run
if the system is up at this preferred time. It is designed to easily run as
a non system application with minimal installation.

## Features

#### 24h run algorithm

Every time `daily-runner` executes the command it tracks the time and if more
than 24 hours have passed since the last run.

#### Preferred run time

A preferred run time can be set via the command line flag `-preferredTime` ie 01:00:00, if
daily-runner is running at this time the command will be executed irrespective
of when the last run was. 

Please note that daily-run does not attempt to run
the backup at exactly the time. It can be up to time adding the interval
setting. For the default 4min interval this means that the command can run from
01:00:00 up to 01:04:10.

#### Configurable check interval

`daily-runner` checks in a loop whether the command needs to run. This loop by
default sleeps 4min between checks to conserve power. This can be adjusted via
command line flag `-interval`.

#### Profiles

### Single process running via pid file locking

## Motivation

I have been using [duply](https://duply.net/index.php/Main_Page) at work with
great success, even used it once to migrate to a new computer, and I wanted to
also setup backups for my home computers. 

At work was running it via cron at
lunchtime but for my home computers which I use randomly & sporadically I
needed to find a different way to run it. I looked at fcron & anacron and I
found them hard to install requiring system changes(fcron an compilation from
sources) plus a steep learning & setup curve to reach my use case. 

Gnome's
[deja dup](https://wiki.gnome.org/Apps/DejaDup) is also an app that tries to
solve the overall problem of home/personal computer backup and has the feature
of "daily", "weekly" running of backup for systems that are not always online.
It's design imposes many defaults, so I was not able to set it up successfully
for my use cases. 
