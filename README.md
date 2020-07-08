# Focus

Focus is a CLI app which edits your `/etc/hosts` file to temporarily stop your
website from connecting to a list of websites which you define.

## Example

```sh
$ focus
Backed up /etc/hosts to ~/.config/focus/hosts.bak
Added entries to /etc/hosts
Blocking websites. Press ctrl+c to stop
```

Focus requires a config file to be stored at `~/.config/focus/focus.yaml`:

```yaml
profiles:
    default:
        - news.ycombinator.com
        - bbc.co.uk
```

You can specify multiple profiles, if there are different websites you want to
block at different times:

```yaml
profiles:
    default:
        - news.ycombinator.com
        - bbc.co.uk
    work:
        - slack.com
        - gmail.com
```

Focus has a built in timer, which will unblock websites and send a desktop
notification once the specified time has elapsed:

```sh
$ focus -timer 25m
```

Valid example durations are `5s`, `5m`, `5h`, which unblock after 5 seconds,
minutes and hours respectively. Negative or zero duration values will block
forever.

## Usage

```
$ focus -h
Usage of focus:
  -profile string
        The profile to use (default "default")
  -timer duration
        Stop blocking after a period of time
```
