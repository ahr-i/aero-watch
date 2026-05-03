import os
import tkinter as tk
from tkinter import ttk

from app.config import load_settings
from app.ui import StreamingApp


def main() -> None:
    root = tk.Tk()
    icon_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)),
        "assets",
        "icon",
        "icon.png",
    )
    if os.path.exists(icon_path):
        icon_image = tk.PhotoImage(file=icon_path)
        root.iconphoto(True, icon_image)
        root._icon_image = icon_image

    style = ttk.Style()
    if "vista" in style.theme_names():
        style.theme_use("vista")

    settings = load_settings()
    StreamingApp(root, settings)
    root.mainloop()


if __name__ == "__main__":
    main()
