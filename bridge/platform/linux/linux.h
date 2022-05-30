#include <gtk/gtk.h>
#include <JavaScriptCore/JavaScript.h>
#include <webkit2/webkit2.h>
#include <libappindicator/app-indicator.h>

extern void go_menu_callback(GtkMenuItem *,int);

extern void go_webview_callback(WebKitUserContentManager *manager, WebKitJavascriptResult *r, int arg);


static void _g_signal_connect(GtkWidget *item, char *action, void *callback, int user) {
  g_signal_connect(item, action, G_CALLBACK(callback), (void *)user);
}

static char *string_from_js_result(WebKitJavascriptResult *r) {
    char *s;
#if WEBKIT_MAJOR_VERSION >= 2 && WEBKIT_MINOR_VERSION >= 22
    JSCValue *value = webkit_javascript_result_get_js_value(r);
    s = jsc_value_to_string(value);
#else
    JSGlobalContextRef ctx = webkit_javascript_result_get_global_context(r);
    JSValueRef value = webkit_javascript_result_get_value(r);
    JSStringRef js = JSValueToStringCopy(ctx, value, NULL);
    size_t n = JSStringGetMaximumUTF8CStringSize(js);
    s = g_new(char, n);
    JSStringGetUTF8CString(js, s, n);
    JSStringRelease(js);
#endif
    return s;
}
