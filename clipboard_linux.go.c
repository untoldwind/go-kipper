#ifndef __CLIPBOARD_GO_C
#define __CLIPBOARD_GO_C

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>

static Atom
_clipboard_atom(Display *display)
{
  return XInternAtom(display, "CLIPBOARD", False);
}

static Atom
_primary_atom()
{
  return XA_PRIMARY;
}

static Atom
_seconary_atom()
{
  return XA_SECONDARY;
}

static Atom
_utf8_string_atom(Display *display)
{
  Atom atom;

  atom = XInternAtom(display, "UTF8_STRING", True);
  return atom == None ? XA_STRING : atom;
}

static Time
_get_timestamp(Display *display, Window window)
{
  XEvent event;

  XChangeProperty (display, window, XA_WM_NAME, XA_STRING, 8,
                   PropModeAppend, NULL, 0);

  while (1) {
    XNextEvent (display, &event);

    if (event.type == PropertyNotify)
      return event.xproperty.time;
  }
}

static char *
_wait_selection(Display *display, Window window, Atom selection, Atom request_target)
{
  XEvent event;
  Atom target;
  int format;
  unsigned long bytesafter, length;
  unsigned char * value, * retval = NULL;
  Bool keep_waiting = True;

  while (keep_waiting) {
    XNextEvent (display, &event);

    switch (event.type) {
    case SelectionNotify:
      if (event.xselection.selection != selection)
        break;

      if (event.xselection.property == None) {
        value = NULL;
        keep_waiting = False;
      } else {
        XGetWindowProperty (event.xselection.display,
                            event.xselection.requestor,
                            event.xselection.property, 0L, 1000000,
                            False, (Atom)AnyPropertyType, &target,
                            &format, &length, &bytesafter, &value);

        if (target != request_target) {
          free(retval);
          retval = NULL;
          keep_waiting = False;
        } else {
          retval = strdup(value);
          XFree(value);
          keep_waiting = False;
        }

        XDeleteProperty (event.xselection.display,
                         event.xselection.requestor,
                         event.xselection.property);
      }
      break;
    default:
      break;
    }
  }

  return retval;
}

static char *
_get_selection(Display *display, Window window, Atom selection, Atom target, Time timestamp)
{
  Atom prop;

  prop = XInternAtom (display, "XSEL_DATA", False);
  XConvertSelection(display, selection, target, prop, window, timestamp);
  XSync(display, False);

  return _wait_selection(display, window, selection, target);
}

static char *
_get_selection_text(Display *display, Window window, Atom selection)
{
  unsigned char *result;
  Time timestamp = _get_timestamp(display, window);

  if ((result = _get_selection(display, window, selection, _utf8_string_atom(display), timestamp)) == NULL) {
    result = _get_selection(display, window, selection, XA_STRING, timestamp);
  }

  return result;
}

#endif
