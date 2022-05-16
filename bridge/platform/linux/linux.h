#include <gtk/gtk.h>
#include <libappindicator/app-indicator.h>

#define true 1
#define false 0

typedef int bool;

bool tray_init();
AppIndicator *tray_indicator_new(char *id, char *png_icon_path, GtkMenuShell *menu);
void tray_poll_events();

GtkMenuShell *menu_new();
void menu_append_menu_item(GtkMenuShell *menu, GtkWidget *item);

GtkWidget *menu_item_new(int id, char *title, bool disabled, bool checked, bool separator);
void menu_item_set_submenu(GtkWidget *parent, GtkWidget *child);

void tray_test();