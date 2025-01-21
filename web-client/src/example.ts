import { WebSocketClient } from './websocket-client';
import { ConnectionState } from './types';

// 创建 WebSocket 客户端实例
const client = new WebSocketClient({
    autoReconnect: true,
    reconnectInterval: 5000,
    heartbeatInterval: 30000,
    debug: true
});

// 设置事件处理器
client.setEventHandlers({
    onMessage: (message: string) => {
        console.log('收到消息:', message);
        try {
            const msgObj = JSON.parse(message);
            if (msgObj.type === 'auth_response') {
                console.log(`认证${msgObj.success ? '成功' : '失败'}: ${msgObj.message || ''}`);
            }
        } catch (e) {
            console.log('原始消息:', message);
        }
    },
    onStateChange: (state: ConnectionState) => {
        const states = ['已断开', '正在连接', '已连接', '重连中'];
        console.log('连接状态:', states[state]);
    },
    onError: (error: Event) => {
        console.error('发生错误:', error);
    },
    onClose: () => {
        console.log('连接已关闭');
    }
});

// 连接到服务器
client.connect(
    'ws://127.0.0.1:9688/ws',
    'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJsb2dpblR5cGUiOiJsb2dpbiIsImxvZ2luSWQiOiJpbV91c2VyOjM5ODEyMDM2NDM0MDE3OCIsInJuU3RyIjoiaUFqQkhQM09RV21ZQklmS0xYeDBBYTZjdmdHb2ZISngiLCJ1c2VySWQiOjM5ODEyMDM2NDM0MDE3OCwidXNlck5hbWUiOiLoqIDpnZnmgKEifQ.jXStT5wSzYvcCHy4PKLGeskhNxNd3Nx-xrVmh0Yaz98'
); 