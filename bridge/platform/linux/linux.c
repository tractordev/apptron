#include "linux.h"

#if 0
static void tray_menu_callback(GtkMenuItem *item, gpointer data) {
  (void)item;
  Tray_Menu_Item *m = (Tray_Menu_Item *)data;
  //m->cb(m);
}

static GtkMenuShell *_build_tray_menu(Tray_Menu_Item *m) {
  GtkMenuShell *menu = (GtkMenuShell *)gtk_menu_new();

  for (; m != NULL && m->text != NULL; m++) {
    GtkWidget *item;

    if (m->separator) {
      item = gtk_separator_menu_item_new();
    } else {
      if (m->submenu != NULL) {
        item = gtk_menu_item_new_with_label(m->text);
        gtk_menu_item_set_submenu(GTK_MENU_ITEM(item), GTK_WIDGET(_build_tray_menu(m->submenu)));
      } else {
        item = gtk_check_menu_item_new_with_label(m->text);
        gtk_check_menu_item_set_active(GTK_CHECK_MENU_ITEM(item), !!m->checked);
      }
      gtk_widget_set_sensitive(item, !m->disabled);
      /*
      if (m->cb != NULL) {
        g_signal_connect(item, "activate", G_CALLBACK(tray_menu_callback), m);
      }
      */
    }

    gtk_widget_show(item);
    gtk_menu_shell_append(menu, item);
  }
  return menu;
}
#endif

bool tray_init()
{
  if (gtk_init_check(0, NULL) == FALSE) {
    return false;
  }

  return true;
}

AppIndicator *tray_indicator_new(char *id, char *png_icon_path, GtkMenuShell *menu)
{
  AppIndicator *result = app_indicator_new(id, "", APP_INDICATOR_CATEGORY_APPLICATION_STATUS);
  
  app_indicator_set_status(result, APP_INDICATOR_STATUS_ACTIVE);

  //app_indicator_set_title(global_app_indicator, title);
  //app_indicator_set_label(global_app_indicator, title, "");

  app_indicator_set_icon_full(result, png_icon_path, "");
  app_indicator_set_menu(result, GTK_MENU(menu));

  return result;
}

void tray_poll_events()
{
  int blocking = 0;
  gtk_main_iteration_do(blocking);
}


GtkMenuShell *menu_new() 
{
  GtkMenuShell *menu = (GtkMenuShell *)gtk_menu_new();
  return menu;
}

void menu_append_menu_item(GtkMenuShell *menu, GtkWidget *item)
{
  gtk_menu_shell_append(menu, item);
}


void menu_item_callback(GtkMenuItem *item, gpointer user_data)
{
  long long int menu_id = (long long int)user_data;

  printf("clicked! %d\n", menu_id);
  fflush(stdout);
}

GtkWidget *menu_item_new(int id, char *title, bool disabled, bool checked, bool separator)
{
  GtkWidget *item = NULL; 

  if (separator)
  {
    item = gtk_separator_menu_item_new();
    gtk_widget_show(item);
  }
  else
  {
    if (checked)
    {
      item = gtk_check_menu_item_new_with_label(title);
      gtk_check_menu_item_set_active(GTK_CHECK_MENU_ITEM(item), !!checked);
    }
    else
    {
      item = gtk_menu_item_new_with_label(title);
    }

    gtk_widget_set_sensitive(item, !disabled);
    gtk_widget_show(item);

    g_signal_connect(item, "activate", G_CALLBACK(menu_item_callback), id);
  }

  return item;
}

void menu_item_set_submenu(GtkWidget *parent, GtkWidget *child)
{
  gtk_menu_item_set_submenu(GTK_MENU_ITEM(parent), GTK_WIDGET(child));
}


void tray_test() {

  tray_init();

  GtkMenuShell *menu = menu_new();
  GtkWidget *item = menu_item_new(1, "Hello", false, false, false);
  menu_append_menu_item(menu, item);

  AppIndicator *indicator = tray_indicator_new("systray", "/home/nick/dev/_projects/apptron/bridge/misc/icon.png", menu);

  //gtk_main_iteration_do(1);

  //tray_poll_events();

  /*
  while (true) {
    gtk_main_iteration_do(1);
  }
  */

}