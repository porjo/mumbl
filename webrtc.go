package main

import (
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

type WebRTCPeer struct {
	pc    *webrtc.RTCPeerConnection
	track *webrtc.RTCTrack
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

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	//	pc.OnICEConnectionStateChange = func(connectionState ice.ConnectionState) {
	//		log.Printf("Connection State has changed %s \n", connectionState.String())
	//	}
	pc.OnICEConnectionStateChange = onStateChange

	// Create a audio track
	opusTrack, err := pc.NewRTCTrack(webrtc.DefaultPayloadTypeOpus, "audio", "mumble")
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
