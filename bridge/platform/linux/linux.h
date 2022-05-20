#include <gtk/gtk.h>
#include <libappindicator/app-indicator.h>

extern void go_menu_callback(GtkMenuItem *,int);

static void _g_signal_connect(GtkWidget *item, char *action, void *callback, int user) {
  g_signal_connect(item, action, G_CALLBACK(callback), (void *)user);
}
