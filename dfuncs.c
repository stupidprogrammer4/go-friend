#include "dfuncs.h"

extern void moveTo(void *data, float to_x, float to_y);
extern void lineTo(void *data, float to_x, float to_y);
extern void quadraticTo(void *data, float control_x, float control_y, float to_x, float to_y);
extern void cubeTo(void *data, float control1_x, float control1_y, float control2_x, float control2_y, float to_x, float to_y);
extern void closePath(void *data);

void move_to(hb_draw_funcs_t *dfuncs, 
             void *draw_data,
             hb_draw_state_t *state,
             float to_x,
             float to_y,
             void *user_data)
{
    moveTo(draw_data, to_x, to_y);
}

void line_to(hb_draw_funcs_t *dfuncs,
             void *draw_data,
             hb_draw_state_t *state,
             float to_x,
             float to_y,
             void *user_data)
{
    lineTo(draw_data, to_x, to_y);
}

void quadratic_to(hb_draw_funcs_t *dfuncs,
                 void *draw_data,
                 hb_draw_state_t *state,
                 float control_x,
                 float control_y,
                 float to_x,
                 float to_y,
                 void *user_data)
{
    quadraticTo(draw_data, control_x, control_y, to_x, to_y);
}

void cube_to(hb_draw_funcs_t *dfuncs,
             void *draw_data,
             hb_draw_state_t *state,
             float control1_x,
             float control1_y,
             float control2_x,
             float control2_y,
             float to_x,
             float to_y,
             void *user_data)
{
    cubeTo(draw_data, control1_x, control1_y, control2_x, control2_y, to_x, to_y);
}

void close_path(hb_draw_funcs_t *dfuncs,
                void *draw_data,
                hb_draw_state_t *state,
                void *user_data)
{
    closePath(draw_data);
}