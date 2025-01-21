import { WebSocketClient } from '../src/websocket-client';
import { ConnectionState } from '../src/types';

declare global {
    interface Window {
        app: App;
    }
}

class App {
    private client: WebSocketClient;
    private messageInput: HTMLInputElement;
    private messagesDiv: HTMLDivElement;
    private statusDiv: HTMLDivElement;

    constructor() {
        // 初始化 WebSocket 客户端
        this.client = new WebSocketClient({
            autoReconnect: true,
            reconnectInterval: 5000,
            heartbeatInterval: 30000,
            debug: true
        });

        // 获取 DOM 元素
        this.messageInput = document.getElementById('message') as HTMLInputElement;
        this.messagesDiv = document.getElementById('messages') as HTMLDivElement;
        this.statusDiv = document.getElementById('status') as HTMLDivElement;

        // 设置事件处理器
        this.setupEventHandlers();
        this.setupUIHandlers();
    }

    private setupEventHandlers(): void {
        this.client.setEventHandlers({
            onMessage: (message: string) => {
                try {
                    const msgObj = JSON.parse(message);
                    let displayMsg = '';
                    let className = '';
                    
                    if (msgObj.type === 'auth_response') {
                        displayMsg = `认证${msgObj.success ? '成功' : '失败'}: ${msgObj.message || ''}`;
                        className = 'auth_response';
                    } else {
                        displayMsg = message;
                    }
                    
                    this.addMessage(displayMsg, className);
                } catch (e) {
                    this.addMessage(message);
                }
            },
            onStateChange: (state: ConnectionState) => {
                const states = ['已断开', '正在连接', '已连接', '重连中'];
                this.statusDiv.textContent = states[state];
                this.statusDiv.className = state === ConnectionState.Connected ? 'connected' : 'disconnected';
            },
            onError: (error: Event) => {
                console.error('发生错误:', error);
                this.addMessage(`错误: ${error}`, 'error');
            },
            onClose: () => {
                console.log('连接已关闭');
                this.addMessage('连接已关闭', 'system');
            }
        });
    }

    private setupUIHandlers(): void {
        // 回车发送消息
        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.sendMessage();
            }
        });
    }

    // 连接到服务器
    connect(): void {
        const url = (document.getElementById('url') as HTMLInputElement).value;
        const token = (document.getElementById('token') as HTMLInputElement).value;
        const deviceType = (document.getElementById('deviceType') as HTMLInputElement).value;

        // 清空消息区域
        this.messagesDiv.innerHTML = '';
        
        this.client.connect(url, token, deviceType);
    }

    // 断开连接
    disconnect(): void {
        this.client.close();
    }

    // 发送消息
    sendMessage(): void {
        const message = this.messageInput.value.trim();
        if (!message) {
            alert('请输入消息内容');
            return;
        }
        
        if (this.client.send(message)) {
            this.addMessage(`发送: ${message}`, 'sent');
            this.messageInput.value = '';
        } else {
            alert('发送失败，请检查连接状态');
        }
    }

    // 添加消息到显示区域
    private addMessage(message: string, className: string = ''): void {
        const div = document.createElement('div');
        div.className = `message ${className}`;
        div.textContent = message;
        this.messagesDiv.appendChild(div);
        this.messagesDiv.scrollTop = this.messagesDiv.scrollHeight;
    }
}

// 初始化应用
window.addEventListener('DOMContentLoaded', () => {
    window.app = new App();
}); 