
//void hello(char *name);
void gomain();

static void invoke(void (*f)()) {
    f();
}

typedef void (*closure)();

void go_main_loop();

typedef enum bool {
    false = 0,
    true = 1,
} bool;

bool window_set_title(int window_id, char *title);

typedef struct Event_Loop {
    unsigned char data[40];
} Event_Loop;

Event_Loop create_event_loop();

int create_window(Event_Loop event_loop);

int run(Event_Loop event_loop, void (*callback)());

/*
#include <stdio.h>

static inline int create_window(Event_Loop event_loop) {
    void *data = (void *)&event_loop;

    printf("[C] create_window\n");
    int size = 40;
    for (int i = 0; i < size; i++) {
      printf("%d, ", ((unsigned char *) data) [i]);
    }

    printf("\n");

    return 0;
};
*/