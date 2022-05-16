#include "linux.h"

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

  if (png_icon_path != NULL)
  {
    app_indicator_set_icon_full(result, png_icon_path, "");
  }

  if (menu != NULL)
  {
    app_indicator_set_menu(result, GTK_MENU(menu));
  }

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
  int menu_id = (int)user_data;

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

void menu_item_set_submenu(GtkWidget *parent, GtkMenuShell *child)
{
  gtk_menu_item_set_submenu(GTK_MENU_ITEM(parent), GTK_WIDGET(child));
}

#if 0
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
#endif