import subprocess
import time
import tkinter as tk
from pathlib import Path
from tkinter import messagebox

import mss
from PIL import Image, ImageTk

from app.ffmpeg import build_stream_command, find_ffmpeg
from app.models import AppSettings, MonitorInfo
from app.monitoring import list_monitors


BG_COLOR = "#F8F9FA"
CARD_COLOR = "#FFFFFF"
SURFACE_COLOR = "#EDEEEF"
SURFACE_LOW = "#F3F4F5"
TEXT_COLOR = "#191C1D"
MUTED_TEXT = "#5F6368"
OUTLINE = "#C1C6D6"
OUTLINE_DARK = "#727785"
PRIMARY = "#1A73E8"
PRIMARY_DARK = "#005BC0"
PRIMARY_SOFT = "#D8E2FF"
ERROR = "#BA1A1A"
SUCCESS = "#1E8E3E"


class StreamingApp:
    def __init__(self, root: tk.Tk, settings: AppSettings):
        self.root = root
        self.root.title("Streaming Test - Ahri Project: aero-watch")
        self.root.geometry("1280x860")
        self.root.configure(bg=BG_COLOR)

        self.settings = settings
        self.ffmpeg_path = find_ffmpeg(settings.ffmpeg_path)
        self.process: subprocess.Popen | None = None
        self.preview_job: str | None = None
        self.metrics_job: str | None = None
        self.preview_images: dict[str, ImageTk.PhotoImage] = {}
        self.monitor_map: dict[str, MonitorInfo] = {}
        self.monitor_cards: dict[str, dict[str, tk.Widget]] = {}
        self.stream_started_at: float | None = None
        self.sct = mss.mss()
        self.topbar_icon: ImageTk.PhotoImage | None = None

        self.selected_monitor = tk.StringVar()
        self.stream_url = tk.StringVar(value=settings.default_rtmp_url)
        self.status_text = tk.StringVar(value="Offline")
        self.monitor_badge_text = tk.StringVar(value="0 Displays Detected")
        self.bitrate_text = tk.StringVar(value="0 Kbps")
        self.fps_text = tk.StringVar(value="0 FPS")
        self.uptime_text = tk.StringVar(value="--:--")
        self.footer_status_text = tk.StringVar(value="Offline")
        self.ffmpeg_text = tk.StringVar(value=self.ffmpeg_path or "FFmpeg not found")

        self.build_ui()
        self.refresh_monitors()
        self.update_status_ui()
        self.root.protocol("WM_DELETE_WINDOW", self.on_close)

    def build_ui(self) -> None:
        self.build_topbar()

        content_shell = tk.Frame(self.root, bg=BG_COLOR)
        content_shell.pack(fill="both", expand=True)

        scrollbar = tk.Scrollbar(content_shell, orient="vertical")
        scrollbar.pack(side="right", fill="y")

        self.content_canvas = tk.Canvas(
            content_shell,
            bg=BG_COLOR,
            highlightthickness=0,
            yscrollcommand=scrollbar.set,
        )
        self.content_canvas.pack(side="left", fill="both", expand=True)
        scrollbar.configure(command=self.content_canvas.yview)

        content = tk.Frame(self.content_canvas, bg=BG_COLOR)
        self.content_window = self.content_canvas.create_window((0, 0), window=content, anchor="nw")
        content.bind("<Configure>", self._update_scrollregion)
        self.content_canvas.bind("<Configure>", self._sync_canvas_width)
        self.bind_scroll_events()

        content_inner = tk.Frame(content, bg=BG_COLOR)
        content_inner.pack(fill="both", expand=True, padx=48, pady=(28, 56))

        welcome = tk.Frame(content_inner, bg=BG_COLOR)
        welcome.pack(fill="x", pady=(0, 24))

        tk.Label(
            welcome,
            text="Ahri Project",
            bg=BG_COLOR,
            fg=TEXT_COLOR,
            font=("Inter", 24, "normal"),
        ).pack(anchor="w")
        tk.Label(
            welcome,
            text="It was created to test the RTMP streaming feature in the Ahri Project. (RTMP 테스트용)",
            bg=BG_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 12, "normal"),
        ).pack(anchor="w", pady=(8, 0))

        self.build_monitor_section(content_inner)
        self.build_destination_section(content_inner)
        self.build_action_section(content_inner)
        self.build_footer()

    def build_topbar(self) -> None:
        topbar = tk.Frame(self.root, bg=CARD_COLOR, height=64, highlightbackground="#E5E7EB", highlightthickness=1)
        topbar.pack(fill="x", side="top")
        topbar.pack_propagate(False)

        left = tk.Frame(topbar, bg=CARD_COLOR)
        left.pack(side="left", padx=20, fill="y")

        icon_path = Path(__file__).resolve().parent.parent / "assets" / "icon" / "icon.png"
        if icon_path.exists():
            icon_image = Image.open(icon_path)
            icon_image.thumbnail((32, 32))
            self.topbar_icon = ImageTk.PhotoImage(icon_image)
            tk.Label(left, image=self.topbar_icon, bg=CARD_COLOR).pack(side="left", pady=16)

        tk.Label(
            left,
            text="Streaming Test",
            bg=CARD_COLOR,
            fg=TEXT_COLOR,
            font=("Inter", 16, "bold"),
        ).pack(side="left", padx=(10, 0))

    def build_monitor_section(self, parent: tk.Widget) -> None:
        section = tk.Frame(parent, bg=BG_COLOR)
        section.pack(fill="x", pady=(0, 24))

        header = tk.Frame(section, bg=BG_COLOR)
        header.pack(fill="x", pady=(0, 12))

        tk.Label(
            header,
            text="Select Monitor",
            bg=BG_COLOR,
            fg=TEXT_COLOR,
            font=("Inter", 16, "bold"),
        ).pack(side="left")
        header_right = tk.Frame(header, bg=BG_COLOR)
        header_right.pack(side="right")

        tk.Button(
            header_right,
            text="Refresh",
            command=self.refresh_monitors,
            bg=BG_COLOR,
            fg=PRIMARY,
            activebackground=PRIMARY_SOFT,
            activeforeground=PRIMARY_DARK,
            relief="flat",
            font=("Inter", 10, "bold"),
            padx=12,
            pady=4,
            cursor="hand2",
        ).pack(side="left", padx=(0, 10))

        tk.Label(
            header_right,
            textvariable=self.monitor_badge_text,
            bg="#DEE0E4",
            fg="#606366",
            font=("Inter", 10, "bold"),
            padx=12,
            pady=4,
        ).pack(side="left")

        self.monitor_grid = tk.Frame(section, bg=BG_COLOR)
        self.monitor_grid.pack(fill="x")
        for column in range(3):
            self.monitor_grid.columnconfigure(column, weight=1)

    def build_destination_section(self, parent: tk.Widget) -> None:
        section = tk.Frame(parent, bg=BG_COLOR)
        section.pack(fill="x", pady=(0, 28))

        tk.Label(
            section,
            text="Stream Destination",
            bg=BG_COLOR,
            fg=TEXT_COLOR,
            font=("Inter", 16, "bold"),
        ).pack(anchor="w", pady=(0, 12))

        card = tk.Frame(
            section,
            bg=CARD_COLOR,
            highlightbackground=OUTLINE,
            highlightthickness=1,
            padx=20,
            pady=18,
        )
        card.pack(fill="x")

        tk.Label(
            card,
            text="Streaming URL (RTMP)",
            bg=CARD_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 10, "bold"),
        ).pack(anchor="w", pady=(0, 8))

        entry_row = tk.Frame(card, bg=CARD_COLOR)
        entry_row.pack(fill="x")

        self.url_entry = tk.Entry(
            entry_row,
            textvariable=self.stream_url,
            bg=BG_COLOR,
            fg=TEXT_COLOR,
            relief="flat",
            highlightbackground=OUTLINE,
            highlightcolor=PRIMARY,
            highlightthickness=1,
            insertbackground=TEXT_COLOR,
            font=("Inter", 11),
        )
        self.url_entry.pack(side="left", fill="x", expand=True, ipady=10)

        tk.Button(
            entry_row,
            text="Copy",
            command=self.copy_stream_url,
            bg=BG_COLOR,
            fg=PRIMARY,
            activebackground=PRIMARY_SOFT,
            activeforeground=PRIMARY_DARK,
            relief="flat",
            font=("Inter", 10, "bold"),
            padx=16,
            pady=10,
            cursor="hand2",
        ).pack(side="left", padx=(12, 0))

        tk.Button(
            entry_row,
            text="Paste",
            command=self.paste_stream_url,
            bg=BG_COLOR,
            fg=PRIMARY,
            activebackground=PRIMARY_SOFT,
            activeforeground=PRIMARY_DARK,
            relief="flat",
            font=("Inter", 10, "bold"),
            padx=16,
            pady=10,
            cursor="hand2",
        ).pack(side="left", padx=(8, 0))

        '''
        tk.Label(
            card,
            text="RTMP 주소 입력 아래 공간",
            bg=CARD_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 10, "normal"),
        ).pack(anchor="w", pady=(12, 0))
        '''

        '''
        tk.Label(
            card,
            textvariable=self.ffmpeg_text,
            bg=CARD_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 9, "normal"),
            wraplength=980,
            justify="left",
        ).pack(anchor="w", pady=(8, 0))
        '''

    def build_action_section(self, parent: tk.Widget) -> None:
        section = tk.Frame(parent, bg=BG_COLOR, highlightbackground=OUTLINE, highlightthickness=0)
        section.pack(fill="x", pady=(0, 16))

        divider = tk.Frame(section, bg=OUTLINE, height=1)
        divider.pack(fill="x", pady=(0, 28))

        buttons = tk.Frame(section, bg=BG_COLOR)
        buttons.pack()

        self.start_button = tk.Button(
            buttons,
            text="Start Streaming",
            command=self.start_stream,
            bg=PRIMARY,
            fg="white",
            activebackground=PRIMARY_DARK,
            activeforeground="white",
            relief="flat",
            font=("Inter", 13, "bold"),
            padx=28,
            pady=14,
            cursor="hand2",
        )
        self.start_button.pack(side="left", padx=8)

        self.stop_button = tk.Button(
            buttons,
            text="Stop Streaming",
            command=self.stop_stream,
            bg=SURFACE_LOW,
            fg=OUTLINE_DARK,
            activebackground=SURFACE_COLOR,
            activeforeground=OUTLINE_DARK,
            relief="flat",
            font=("Inter", 13, "bold"),
            padx=28,
            pady=14,
            state="disabled",
            cursor="hand2",
        )
        self.stop_button.pack(side="left", padx=8)

        tk.Label(
            section,
            textvariable=self.status_text,
            bg=BG_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 10, "normal"),
        ).pack(pady=(14, 0))

        metrics = tk.Frame(section, bg=BG_COLOR)
        metrics.pack(pady=(24, 0))

        self.build_metric(metrics, self.bitrate_text, "Bitrate").pack(side="left", padx=28)
        tk.Frame(metrics, bg=OUTLINE, width=1, height=42).pack(side="left")
        self.build_metric(metrics, self.fps_text, "Frame Rate").pack(side="left", padx=28)
        tk.Frame(metrics, bg=OUTLINE, width=1, height=42).pack(side="left")
        self.build_metric(metrics, self.uptime_text, "Uptime").pack(side="left", padx=28)

    def build_metric(self, parent: tk.Widget, value_var: tk.StringVar, label: str) -> tk.Frame:
        frame = tk.Frame(parent, bg=BG_COLOR)
        tk.Label(
            frame,
            textvariable=value_var,
            bg=BG_COLOR,
            fg=PRIMARY,
            font=("Inter", 18, "bold"),
        ).pack()
        tk.Label(
            frame,
            text=label,
            bg=BG_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 10, "normal"),
        ).pack(pady=(4, 0))
        return frame

    def build_footer(self) -> None:
        footer = tk.Frame(self.root, bg=CARD_COLOR, height=40, highlightbackground="#E5E7EB", highlightthickness=1)
        footer.pack(fill="x", side="bottom")
        footer.pack_propagate(False)

        left = tk.Frame(footer, bg=CARD_COLOR)
        left.pack(side="left", padx=20, fill="y")

        self.status_dot = tk.Canvas(left, width=12, height=12, bg=CARD_COLOR, highlightthickness=0)
        self.status_dot.pack(side="left", pady=12)
        self.status_dot_id = self.status_dot.create_oval(2, 2, 10, 10, fill=ERROR, outline=ERROR)

        tk.Label(
            left,
            textvariable=self.footer_status_text,
            bg=CARD_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 10, "normal"),
        ).pack(side="left", padx=(6, 12))

        tk.Label(
            left,
            text="v1.2.0-test",
            bg=CARD_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 10, "normal"),
        ).pack(side="left")

        right = tk.Frame(footer, bg=CARD_COLOR)
        right.pack(side="right", padx=20, fill="y")

        tk.Label(
            right,
            text="CPU: --",
            bg=CARD_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 10, "normal"),
        ).pack(side="left", padx=10, pady=10)
        tk.Label(
            right,
            text="Latency: 24ms",
            bg=CARD_COLOR,
            fg=MUTED_TEXT,
            font=("Inter", 10, "normal"),
        ).pack(side="left", padx=10, pady=10)

    def refresh_monitors(self) -> None:
        monitors = list_monitors()
        self.monitor_map = {monitor.label: monitor for monitor in monitors}
        labels = list(self.monitor_map)
        self.monitor_badge_text.set(f"{len(labels)} Displays Detected")

        if not labels:
            self.selected_monitor.set("")
            self.clear_monitor_cards()
            self.status_text.set("No monitors detected.")
            self.update_status_ui()
            return

        if self.selected_monitor.get() not in self.monitor_map:
            self.selected_monitor.set(labels[0])

        self.render_monitor_cards(labels)
        self.update_previews()

    def clear_monitor_cards(self) -> None:
        for child in self.monitor_grid.winfo_children():
            child.destroy()
        self.monitor_cards.clear()
        self.preview_images.clear()

    def render_monitor_cards(self, labels: list[str]) -> None:
        self.clear_monitor_cards()

        for index, label in enumerate(labels):
            monitor = self.monitor_map[label]
            card = tk.Frame(
                self.monitor_grid,
                bg=CARD_COLOR,
                highlightbackground=OUTLINE,
                highlightthickness=1,
                bd=0,
                padx=14,
                pady=14,
            )
            card.grid(row=index // 3, column=index % 3, padx=10, pady=10, sticky="nsew")
            card.bind("<Button-1>", lambda _event, selected=label: self.select_monitor(selected))

            preview_shell = tk.Frame(
                card,
                bg=SURFACE_COLOR,
                height=196,
            )
            preview_shell.pack(fill="x")
            preview_shell.pack_propagate(False)
            preview_shell.bind("<Button-1>", lambda _event, selected=label: self.select_monitor(selected))

            preview = tk.Label(
                preview_shell,
                text="Loading preview...",
                bg=SURFACE_COLOR,
                fg=MUTED_TEXT,
                font=("Inter", 10),
            )
            preview.pack(fill="both", expand=True, padx=8, pady=8)
            preview.bind("<Button-1>", lambda _event, selected=label: self.select_monitor(selected))

            info = tk.Frame(card, bg=CARD_COLOR)
            info.pack(fill="x", pady=(12, 0))
            info.bind("<Button-1>", lambda _event, selected=label: self.select_monitor(selected))

            title = tk.Label(
                info,
                text=self.monitor_title(index, monitor),
                bg=CARD_COLOR,
                fg=TEXT_COLOR,
                font=("Inter", 11, "bold"),
                anchor="w",
            )
            title.pack(anchor="w")
            title.bind("<Button-1>", lambda _event, selected=label: self.select_monitor(selected))

            subtitle = tk.Label(
                info,
                text=self.monitor_subtitle(monitor),
                bg=CARD_COLOR,
                fg=MUTED_TEXT,
                font=("Inter", 9),
                anchor="w",
            )
            subtitle.pack(anchor="w", pady=(4, 10))
            subtitle.bind("<Button-1>", lambda _event, selected=label: self.select_monitor(selected))

            action = tk.Button(
                info,
                text="Select",
                command=lambda selected=label: self.select_monitor(selected),
                bg=CARD_COLOR,
                fg=PRIMARY,
                activebackground=PRIMARY_SOFT,
                activeforeground=PRIMARY_DARK,
                relief="flat",
                font=("Inter", 10, "bold"),
                padx=14,
                pady=6,
                cursor="hand2",
            )
            action.pack(anchor="e")

            self.monitor_cards[label] = {
                "card": card,
                "preview": preview,
                "title": title,
                "subtitle": subtitle,
                "action": action,
            }

        self.update_monitor_selection_styles()

    def monitor_title(self, index: int, monitor: MonitorInfo) -> str:
        size_label = self.resolution_label(monitor)
        return f"Monitor {index + 1} ({size_label})"

    def monitor_subtitle(self, monitor: MonitorInfo) -> str:
        orientation = "Landscape" if monitor.width >= monitor.height else "Portrait"
        primary_label = " • Primary" if monitor.x == 0 and monitor.y == 0 else ""
        return f"{monitor.width} x {monitor.height} • {orientation}{primary_label}"

    def resolution_label(self, monitor: MonitorInfo) -> str:
        if monitor.width >= 3840:
            return "4K"
        if monitor.width >= 2560:
            return "UHD"
        if monitor.width >= 1920:
            return "HD"
        return "Display"

    def select_monitor(self, label: str) -> None:
        self.selected_monitor.set(label)
        self.update_monitor_selection_styles()

    def update_monitor_selection_styles(self) -> None:
        selected = self.selected_monitor.get()
        for label, widgets in self.monitor_cards.items():
            is_selected = label == selected
            widgets["card"].configure(
                highlightbackground=PRIMARY if is_selected else OUTLINE,
                highlightthickness=2 if is_selected else 1,
                bg=CARD_COLOR,
            )
            widgets["action"].configure(
                text="Selected" if is_selected else "Select",
                bg=PRIMARY if is_selected else CARD_COLOR,
                fg="white" if is_selected else PRIMARY,
                activebackground=PRIMARY_DARK if is_selected else PRIMARY_SOFT,
                activeforeground="white" if is_selected else PRIMARY_DARK,
            )

    def update_previews(self) -> None:
        if self.preview_job:
            self.root.after_cancel(self.preview_job)
            self.preview_job = None

        for label, monitor in self.monitor_map.items():
            widgets = self.monitor_cards.get(label)
            if not widgets:
                continue

            try:
                shot = self.sct.grab(
                    {
                        "left": monitor.x,
                        "top": monitor.y,
                        "width": monitor.width,
                        "height": monitor.height,
                    }
                )
                image = Image.frombytes("RGB", shot.size, shot.rgb)
                image.thumbnail((304, 180))
                photo = ImageTk.PhotoImage(image)
                self.preview_images[label] = photo
                widgets["preview"].configure(image=photo, text="")
            except Exception as exc:
                widgets["preview"].configure(text=f"Preview unavailable\n{exc}", image="")

        self.preview_job = self.root.after(self.settings.preview_interval_ms, self.update_previews)

    def copy_stream_url(self) -> None:
        self.root.clipboard_clear()
        self.root.clipboard_append(self.stream_url.get().strip())
        self.status_text.set("Streaming URL copied to clipboard.")
        self.update_status_ui()

    def paste_stream_url(self) -> None:
        try:
            clipboard_text = self.root.clipboard_get().strip()
        except tk.TclError:
            self.status_text.set("Clipboard is empty.")
            self.update_status_ui()
            return

        self.stream_url.set(clipboard_text)
        self.status_text.set("Streaming URL pasted from clipboard.")
        self.update_status_ui()

    def start_stream(self) -> None:
        monitor = self.monitor_map.get(self.selected_monitor.get())
        url = self.stream_url.get().strip()

        if not monitor:
            messagebox.showerror("Missing monitor", "Select a monitor before starting the stream.")
            return
        if not url.startswith("rtmp://"):
            messagebox.showerror("Invalid RTMP URL", "Enter an RTMP URL that starts with rtmp://")
            return
        if not self.ffmpeg_path:
            messagebox.showerror(
                "FFmpeg not found",
                "FFmpeg could not be found. Install it or set AEROWATCH_FFMPEG.",
            )
            return
        if self.process and self.process.poll() is None:
            messagebox.showinfo("Streaming", "A stream is already running.")
            return

        command = build_stream_command(self.ffmpeg_path, monitor, url)

        try:
            self.process = subprocess.Popen(
                command,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )
        except Exception as exc:
            messagebox.showerror("Stream start failed", str(exc))
            return

        self.stream_started_at = time.time()
        self.status_text.set(f"Streaming {monitor.label} to {url}")
        self.bitrate_text.set("6000 Kbps")
        self.fps_text.set("30 FPS")
        self.start_button.configure(state="disabled", bg=PRIMARY_DARK)
        self.stop_button.configure(state="normal", bg=SURFACE_COLOR, fg=TEXT_COLOR)
        self.update_status_ui()
        self.root.after(1000, self.check_process)
        self.update_stream_metrics()

    def update_stream_metrics(self) -> None:
        if self.metrics_job:
            self.root.after_cancel(self.metrics_job)
            self.metrics_job = None

        if not self.process or self.process.poll() is not None or not self.stream_started_at:
            return

        elapsed = int(time.time() - self.stream_started_at)
        minutes, seconds = divmod(elapsed, 60)
        hours, minutes = divmod(minutes, 60)
        if hours:
            self.uptime_text.set(f"{hours:02d}:{minutes:02d}:{seconds:02d}")
        else:
            self.uptime_text.set(f"{minutes:02d}:{seconds:02d}")

        self.metrics_job = self.root.after(1000, self.update_stream_metrics)

    def check_process(self) -> None:
        if not self.process:
            return

        return_code = self.process.poll()
        if return_code is None:
            self.root.after(1000, self.check_process)
            return

        self.process = None
        self.stream_started_at = None
        self.status_text.set(f"Streaming stopped unexpectedly (exit code {return_code}).")
        self.reset_stream_metrics()
        self.start_button.configure(state="normal", bg=PRIMARY)
        self.stop_button.configure(state="disabled", bg=SURFACE_LOW, fg=OUTLINE_DARK)
        self.update_status_ui()

    def stop_stream(self) -> None:
        if self.process and self.process.poll() is None:
            self.process.terminate()
            try:
                self.process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                self.process.kill()

        self.process = None
        self.stream_started_at = None
        self.status_text.set("Stream stopped.")
        self.reset_stream_metrics()
        self.start_button.configure(state="normal", bg=PRIMARY)
        self.stop_button.configure(state="disabled", bg=SURFACE_LOW, fg=OUTLINE_DARK)
        self.update_status_ui()

    def reset_stream_metrics(self) -> None:
        if self.metrics_job:
            self.root.after_cancel(self.metrics_job)
            self.metrics_job = None
        self.bitrate_text.set("0 Kbps")
        self.fps_text.set("0 FPS")
        self.uptime_text.set("--:--")

    def update_status_ui(self) -> None:
        is_live = self.process is not None and self.process.poll() is None
        self.footer_status_text.set("Live" if is_live else "Offline")
        self.status_dot.itemconfig(self.status_dot_id, fill=SUCCESS if is_live else ERROR, outline=SUCCESS if is_live else ERROR)

    def bind_scroll_events(self) -> None:
        self.content_canvas.bind_all("<MouseWheel>", self._on_mousewheel)
        self.content_canvas.bind_all("<Button-4>", self._on_mousewheel)
        self.content_canvas.bind_all("<Button-5>", self._on_mousewheel)

    def _update_scrollregion(self, _event: tk.Event) -> None:
        self.content_canvas.configure(scrollregion=self.content_canvas.bbox("all"))

    def _sync_canvas_width(self, event: tk.Event) -> None:
        self.content_canvas.itemconfigure(self.content_window, width=event.width)

    def _on_mousewheel(self, event: tk.Event) -> None:
        if getattr(event, "num", None) == 4:
            self.content_canvas.yview_scroll(-1, "units")
            return
        if getattr(event, "num", None) == 5:
            self.content_canvas.yview_scroll(1, "units")
            return
        delta = int(-1 * (event.delta / 120))
        self.content_canvas.yview_scroll(delta, "units")

    def on_close(self) -> None:
        if self.preview_job:
            self.root.after_cancel(self.preview_job)
        if self.metrics_job:
            self.root.after_cancel(self.metrics_job)
        self.content_canvas.unbind_all("<MouseWheel>")
        self.content_canvas.unbind_all("<Button-4>")
        self.content_canvas.unbind_all("<Button-5>")
        self.stop_stream()
        self.sct.close()
        self.root.destroy()
