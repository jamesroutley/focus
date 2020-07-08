# Focus

Focus is a CLI app which edits your `/etc/hosts` file to temporarily stop your
website from connecting to a list of websites which you define. It supports
multiple lists of websites to block (we call them 'profiles'), and has a
built-in timer which will unblock and send you a desktop notification after the
specified time has elapsed.

⚠️ Focus requires `sudo` to work, and modifies and deletes system files - if
you want to use it please do so at your own risk.

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

### Profiles

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

### Timer

```sh
$ focus -timer 25m
Backed up /etc/hosts to ~/.config/focus/hosts.bak
Added entries to /etc/hosts
Blocking websites. Press ctrl+c to stop
25m remaining
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
