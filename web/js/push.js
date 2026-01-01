"use strict";

// 开始推流按键
const btnStart = document.getElementById('StartPush');
// 停止推流按键
const btnStop = document.getElementById('StopPush');
// 启用音频选项
const checkAudio = document.getElementById('enableAudio');
// 启用音频选项
const checkVideo = document.getElementById('enableVideo');
// 启用共享屏幕选项
const checkScreen = document.getElementById('enableScreen');
// 本端视频
const localVideo = document.getElementById('localVideo');
// 推流url
const pushUrl = document.getElementById('url');

// webrtc peer connection
var peerConn = null;

btnStart.addEventListener('click', async function() {
	try {
        if (btnStart.disabled) {
            alert('请先停止推流');
            return;
        }

        // 修改按键状态
        btnStart.disabled = true;
        btnStart.style.backgroundColor = '#e74c3c';
        btnStart.style.color = 'white';

        btnStop.disabled = false;
        btnStop.style.backgroundColor = '#2ecc71';
        btnStop.style.color = 'white';

        // 获取媒体流选项
        const audioEnabled = checkAudio.checked;
        checkAudio.disabled = true;

        const videoEnabled = checkVideo.checked;
        checkVideo.disabled = true;

        const screenEnabled = checkScreen.checked;
        checkScreen.disabled = true;

        // 获取本地媒体流
        var localStream = null;
        if (screenEnabled) {
            localStream = await navigator.mediaDevices.getDisplayMedia({
                video: videoEnabled,
                audio: audioEnabled
            });
        }
        else {
            localStream = await navigator.mediaDevices.getUserMedia({
                video: videoEnabled,
                audio: audioEnabled
            });
        }
        
        // 将流显示在video元素中
        localVideo.srcObject = localStream;

        // 启动推流
        startPush(localStream, pushUrl.value);
    } 
	catch (error) {
        console.error('Start push error: ', error);
        reset();
    }
});

function reset() {
    btnStart.disabled = false;
    btnStart.style.backgroundColor = '#2ecc71';
    btnStart.style.color = 'white';

    btnStop.disabled = true;
    btnStop.style.backgroundColor = '#e74c3c';
    btnStop.style.color = 'white';

    checkAudio.disabled  = false;
    checkVideo.disabled  = false;
    checkScreen.disabled = false;

    localVideo.srcObject = null;

    if (null != peerConn) {
        peerConn.close();
        peerConn = null;
    }
}

function startPush(localStream, url) {
    peerConn = new RTCPeerConnection();
    peerConn.addStream(localStream);

    // 设置ice事件
    peerConn.oniceconnectionstatechange = function(event) {
        console.log('On ice state change event: ', event, '. state: ', peerConn.iceConnectionState);
    }

    peerConn.onicecandidate = function(event) {
        if (null != event.candidate) {
            console.log('On ice candidate: ', event.candidate.candidate);
            peerConn.addIceCandidate(event.candidate);
        }
    }

    // 创建sdp offer
    peerConn.createOffer()
        .then(function(desc) {
            // 设置本端sdp
            console.log('Local sdp: ', desc.sdp);
            peerConn.setLocalDescription(desc);

            const reqBody = {
                type: desc.type,
                sdp: desc.sdp
            };

            // 发送sdp offer给对端
            fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(reqBody)
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`response status: ${response.status}`);
                }

                console.log("Receive response: ", response);
                const resp = response.json();
            })
            .catch(error => {
                console.error('Start push error: ', error);
                reset();
            });
        })
        .catch(function(error) {
            console.error('Start push error: ', error);
            reset();
        });
}