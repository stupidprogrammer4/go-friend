#include <harfbuzz/hb.h>

void move_to(hb_draw_funcs_t *dfuncs, 
             void *draw_data,
             hb_draw_state_t *state,
             float to_x,
             float to_y,
             void *user_data);

void line_to(hb_draw_funcs_t *dfuncs,
             void *draw_data,
             hb_draw_state_t *state,
             float to_x,
             float to_y,
             void *user_data);

void quadratic_to(hb_draw_funcs_t *dfuncs,
                 void *draw_data,
                 hb_draw_state_t *state,
                 float control_x,
                 float control_y,
                 float to_x,
                 float to_y,
                 void *user_data);

void cube_to(hb_draw_funcs_t *dfuncs,
             void *draw_data,
             hb_draw_state_t *state,
             float control1_x,
             float control1_y,
             float control2_x,
             float control2_y,
             float to_x,
             float to_y,
             void *user_data);

void close_path(hb_draw_funcs_t *dfuncs,
                void *draw_data,
                hb_draw_state_t *state,
                void *user_data);