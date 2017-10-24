# i3-battery-nagbar

Shows i3 nag bar when battery discharging threshold had been achieved, removes
nag bar when battery charge is present.

Stupid and simple program without magic.

# Installation

Happy Arch Linux users shall find package in AUR, other guys should install it
in go get manner.

# Options

- `--threshold <int>` - Threshold when i3-nagbar should appear, default: `15`.
- `--message <msg>` - Message for i3-nagbar,
    default: `Too low charge of battery: {{ .percentage }}`.
- `--interval <duration>` - Interval for checking battery state, default: `1s`'.
- `-h --help` - Show this message.
- `--version` - Show version.

# License

MIT.
