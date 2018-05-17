package worker

var Channel = make(chan func(), 100)

func CallSynchronized(fn func()) {
	Channel <- fn
}
