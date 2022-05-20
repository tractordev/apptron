# Linux Native Package

## Install

Before building on Linux, make sure you have the following installed:

```bash
sudo apt install build-essential
sudo apt install libx11-dev
sudo apt install libgtk-3-dev
sudo apt install libappindicator3-dev
```

NOTE(nick): Most of the windowing stuff _only_ needs X11 (and could even dynamically link against it!), but as far as I can tell there isn't an easy way to get a system tray icon and menu without at least GTK (libappindicator adds a simplified API on top of that).