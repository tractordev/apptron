
void hello(char *name);
void gomain();

static void invoke(void (*f)()) {
    f();
}

typedef void (*closure)();

void go_main_loop();

int run(void (*user_callback)());
