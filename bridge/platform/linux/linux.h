#include <gtk/gtk.h>
#include <libappindicator/app-indicator.h>

#define true 1
#define false 0

typedef int bool;

typedef void (*Menu_Callback)(int);
typedef void (*Closure)();

bool tray_init();
AppIndicator *tray_indicator_new(char *id, char *png_icon_path, GtkMenuShell *menu);
void tray_poll_events();

GtkMenuShell *menu_new();
void menu_append_menu_item(GtkMenuShell *menu, GtkWidget *item);
void menu_set_callback(Menu_Callback  callback);

GtkWidget *menu_item_new(int id, char *title, bool disabled, bool checked, bool separator);
void menu_item_set_submenu(GtkWidget *parent, GtkMenuShell *child);

void tray_test();