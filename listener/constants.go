package listener

const (
	ESCAPE        rune = 27
	CSI           rune = '['
	ENTER         rune = 10
	BACKSPACE     rune = 127
	TAB           rune = 9
	KEY_UP        rune = -2
	KEY_DOWN      rune = -3
	KEY_LEFT      rune = -4
	KEY_RIGHT     rune = -5
	KEY_INSERT    rune = -6
	KEY_HOME      rune = -7
	KEY_PAGE_UP   rune = -8
	KEY_PAGE_DOWN rune = -9
	KEY_DELETE    rune = -10
	KEY_END       rune = -11
	KEY_F1        rune = -12
	KEY_F2        rune = -13
	KEY_F3        rune = -14
	KEY_F4        rune = -15
	KEY_F5        rune = -16
	KEY_F6        rune = -17
	KEY_F7        rune = -18
	KEY_F8        rune = -19
	KEY_F9        rune = -20
	KEY_F10       rune = -21
	KEY_F11       rune = -22
	KEY_F12       rune = -23
)

type eventType int

const (
	KeyPressed eventType = iota
	MouseMove
	MouseLeftClick
	MouseLeftRelease
	MouseRightClick
	MouseRightRelease
	MouseMiddleClick
	MouseMiddleRelease
	ScrollUp
	ScrollDown
	Unknown
)
