// +build mobile

/*
Copyright (C) 2016 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package platform

const maxEvents = 128

var InputEventChan = make(chan interface{}, maxEvents)

func Init() error {
	return nil
}

func Shutdown() {
}

func PollEvent() Event {
	for {
		select {
		case _, ok := <-InputEventChan:
			if ok {
				return &MouseButtonEvent{X: 0, Y: 0, Button: 0, Type: MouseButtonDown}
			} else {
				return QuitEvent{}
			}
		default:
			return nil
		}
	}
}