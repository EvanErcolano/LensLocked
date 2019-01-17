package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlDanger  = "success"

	// AlertMsgGeneric is displayed when any random error occurs
	AlertMsgGeneric = "Something with wrong. Please try again."
)

// Alert is used to render Bootstrap Alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data to come in
type Data struct {
	Alert *Alert
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	// if the error impls the publicErr interace
	// the okay variable will be set to true
	// and then pErr will be set to the variable cast to public err type
	// else will be a generic error
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}

}

// PublicError
type PublicError interface {
	error
	Public() string
}
