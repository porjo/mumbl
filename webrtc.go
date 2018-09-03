/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

type WebRTCPeer struct {
	pc    *webrtc.RTCPeerConnection
	track *webrtc.RTCTrack
}

func (w *WebRTCPeer) Close() error {
	return w.pc.Close()
}

func NewPC(offerSd string, onStateChange func(connectionState ice.ConnectionState)) (*WebRTCPeer, error) {
	// Setup the codecs you want to use.
	// We'll use the default ones but you can also define your own
	webrtc.RegisterDefaultCodecs()

	// Create a new RTCPeerConnection
	pc, err := webrtc.New(webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	pc.OnICEConnectionStateChange = onStateChange

	// Create a audio track
	opusTrack, err := pc.NewRTCTrack(webrtc.DefaultPayloadTypeOpus, "audio", "mumblet")
	if err != nil {
		return nil, err
	}
	_, err = pc.AddTrack(opusTrack)
	if err != nil {
		return nil, err
	}

	// Set the remote SessionDescription
	offer := webrtc.RTCSessionDescription{
		Type: webrtc.RTCSdpTypeOffer,
		Sdp:  offerSd,
	}
	if err := pc.SetRemoteDescription(offer); err != nil {
		return nil, err
	}

	peer := &WebRTCPeer{pc: pc, track: opusTrack}

	return peer, nil
}
