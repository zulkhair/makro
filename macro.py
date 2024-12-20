import win32gui
import win32api
import win32con
import time

def send_key_to_window(window_title, key_code):
    """
    Sends a single key press (down and up) to a window without focusing it.

    Args:
        window_title (str): The title of the window to send the key to.
        key_code (int): The virtual key code of the key to send.
    """
    # Find the window handle by title
    hwnd = win32gui.FindWindow(None, window_title)
    if not hwnd:
        print(f"Window with title '{window_title}' not found.")
        return False  # Indicate failure to find window
    
    # Send the key press (down and up)
    win32api.SendMessage(hwnd, win32con.WM_KEYDOWN, key_code, 0)
    win32api.SendMessage(hwnd, win32con.WM_KEYUP, key_code, 0)
    return True  # Indicate success

# Key codes: Tab=0x09, Space=0x20
window_title = "RF Online"

try:
    while True:  # Infinite loop
        # Send Tab key with a 50-millisecond delay
        if not send_key_to_window(window_title, 0x09):
            break  # Exit loop if the window is not found
        time.sleep(0.05)  # 50 ms delay
        
        # Send Space key with a 3-second delay
        if not send_key_to_window(window_title, 0x20):
            break  # Exit loop if the window is not found
        time.sleep(3.0)  # 3 s delay

except KeyboardInterrupt:
    print("\nProgram stopped by user.")
