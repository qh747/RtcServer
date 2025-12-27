'use strict';

/** ------------------------------------- 界面按键 ----------------------------------------- */

// 开始推流按键
const btnStartPush = document.getElementById('startPushButton');
// 停止推流按键
const btnStopPush = document.getElementById('stopPushButton');

// 开始拉流按键
const btnStartPull = document.getElementById('startPullButton');
// 停止拉流按键
const btnStopPull = document.getElementById('stopPullButton');

// 推流视频
const pushVideoScreen = document.getElementById('localVideo');
// 拉流视频
const pullVideoScreen = document.getElementById('remoteVideo');
 
/** ------------------------------------- 全局变量 ----------------------------------------- */

// 推流端peerConnection
var pushPeerConn = null;

// 推流端媒体流
var pushMediaStream = null;

// 拉流端peerConnection
var pullPeerConn = null;

/** ------------------------------------- 界面函数 ----------------------------------------- */

// 推流按钮点击事件处理函数
btnStartPush.addEventListener('click', async function() {
	try {
        // 使用getDisplayMedia API请求屏幕共享
        const localStream = await navigator.mediaDevices.getDisplayMedia({
            video: { cursor: "always" }, // 显示鼠标指针
            audio: false                 // 不需要音频
        });
        
        // 将流显示在video元素中
        pushVideoScreen.srcObject = localStream;
        
		// 设置按键状态为true, 防止重复点击
        btnStartPush.disabled = true;
        btnStopPush.disabled = false;

        // 调用本端共享屏幕打开成功回调函数
        onLocalScreenOpenSuccess(localStream, {
            // 接收对端音频: false
            offerToReceiveAudio: false,

            // 接收对端视频: false
            offerToReceiveVideo: false
        });
    } 
	catch (error) {
        console.error('Share screen error: ', error);
    }
});

// 停止推流按钮点击事件处理函数
btnStopPush.addEventListener('click', async function() {
	try {
        // 关闭video元素中的流显示
        pushVideoScreen.srcObject = null;
        
		// 设置按键状态为true, 防止重复点击
        btnStartPush.disabled = false;
        btnStopPush.disabled = true;

        // 关闭本端peer connection
        if (null != pushPeerConn) {
            pushPeerConn.close();
            pushPeerConn = null;
        }
    } 
	catch (error) {
        console.error('Stop share screen error: ', error);
    }
});

// 拉流按钮点击事件处理函数
btnStartPull.addEventListener('click', async function() {
	try {
        if (false === btnStartPush.disabled) {
            alert('请先开始推流');
        }

        // 将流显示在video元素中
        if (null != pushMediaStream) {
            pullVideoScreen.srcObject = pushMediaStream;
        }
        
		// 设置按键状态为true, 防止重复点击
        btnStartPull.disabled = true;
        btnStopPull.disabled = false;
    } 
	catch (error) {
        console.error('Pull share screen error: ', error);
    }
});

// 停止拉流按钮点击事件处理函数
btnStopPull.addEventListener('click', async function() {
	try {
        // 关闭video元素中的流显示
        pullVideoScreen.srcObject = null;
        
		// 设置按键状态为true, 防止重复点击
        btnStartPull.disabled = false;
        btnStopPull.disabled = true;

        // 关闭本端peer connection
        if (null != pullPeerConn) {
            pullPeerConn.close();
            pullPeerConn = null;
        }

        // 或者触发停止推流按钮点击事件
        btnStopPush.dispatchEvent(new Event('click'));
    } 
	catch (error) {
        console.error('Stop share screen error: ', error);
    }
});

/** ------------------------------------- 功能函数 ----------------------------------------- */

// 本端共享屏幕打开成功回调函数
function onLocalScreenOpenSuccess(stream, option) {
    // 推流peerConnection
    {
        // 添加媒体流到本端peerConnection
        pushPeerConn = new RTCPeerConnection();
        pushPeerConn.addStream(stream);

        // 设置推流ice
        pushPeerConn.oniceconnectionstatechange = function(event) {
            console.log('Push peer connection state change: ', pushPeerConn.iceConnectionState);
        }

        pushPeerConn.onicecandidate = function(event) {
            if (null != event.candidate) {
                console.log('Push peer connection ice candidate: ', event.candidate.candidate);
                pullPeerConn.addIceCandidate(event.candidate);
            }
        }

        // 创建sdp offer
        pushPeerConn.createOffer(option)
            .then(function(description) {
                // 设置本端sdp
                console.log('Push peer connection create offer success: ', description.sdp);
                pushPeerConn.setLocalDescription(description);

                // 发送sdp offer给对端
                pullPeerConn.setRemoteDescription(description);
                pullPeerConn.createAnswer()
                    .then(function(description) { 
                        console.log('Pull peer connection create answer success: ', description.sdp);
                        pullPeerConn.setLocalDescription(description);
                        pushPeerConn.setRemoteDescription(description);
                    })
                    .catch(function(error) {
                        console.error('Pull peer connection create answer error: ', error);
                    });
            })
            .catch(function(error) {
                console.error('Push peer connection create offer error: ', error);
            });
    }

    // 拉流peerConnection
    {
        pullPeerConn = new RTCPeerConnection();

        // 设置推流ice
        pullPeerConn.oniceconnectionstatechange = function(event) {
            console.log('Pull peer connection state change: ', pullPeerConn.iceConnectionState);
        }

        pullPeerConn.onicecandidate = function(event) {
            if (null != event.candidate) {
                console.log('Pull peer connection ice candidate: ', event.candidate.candidate);
                pushPeerConn.addIceCandidate(event.candidate);
            }
        }

        pullPeerConn.onaddstream = function(event) {
            // 显示对端视频
            pushMediaStream = event.stream;
        }
    }
}