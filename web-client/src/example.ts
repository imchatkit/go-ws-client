import { WebSocketClient } from './websocket-client';
import { ConnectionState } from './types';

class WebSocketDemo {
    private client: WebSocketClient | null = null;
    private messageContainer!: HTMLElement;
    private statusElement!: HTMLElement;

    constructor() {
        this.initializeElements();
        this.bindEvents();
    }

    private initializeElements() {
        this.messageContainer = document.getElementById('messages') as HTMLElement;
        this.statusElement = document.getElementById('status') as HTMLElement;
        
        // 添加回车发送功能
        const messageInput = document.getElementById('messageInput') as HTMLInputElement;
        messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.sendMessage();
            }
        });
    }

    private bindEvents() {
        document.getElementById('connectBtn')?.addEventListener('click', () => this.connect());
        document.getElementById('disconnectBtn')?.addEventListener('click', () => this.disconnect());
        document.getElementById('sendBtn')?.addEventListener('click', () => this.sendMessage());
    }

    private updateStatus(state: ConnectionState) {
        const states = ["已断开", "正在连接", "已连接", "重连中"];
        this.statusElement.textContent = states[state] || "未知状态";
        this.statusElement.className = state === ConnectionState.Connected ? 'connected' : 'disconnected';
    }

    private addMessage(message: string, type: string = '') {
        try {
            // 尝试解析 JSON
            const msgObj = JSON.parse(message);
            let displayMsg = '';
            
            // 根据消息类型显示不同的内容
            if (msgObj.type === 'auth_response') {
                displayMsg = `认证${msgObj.success ? '成功' : '失败'}: ${msgObj.message || ''}`;
            } else {
                displayMsg = message;
            }
            
            this.messageContainer.innerHTML += `<div class="message ${type} ${msgObj.type}">${displayMsg}</div>`;
        } catch (e) {
            // 如果不是 JSON，直接显示原始消息
            this.messageContainer.innerHTML += `<div class="message ${type}">${message}</div>`;
        }
        
        // 滚动到底部
        this.messageContainer.scrollTop = this.messageContainer.scrollHeight;
    }

    public connect() {
        const url = (document.getElementById('url') as HTMLInputElement).value;
        const token = (document.getElementById('token') as HTMLInputElement).value;
        const deviceType = (document.getElementById('deviceType') as HTMLInputElement).value;

        // 清空消息区域
        this.messageContainer.innerHTML = '';

        // 创建新的客户端实例
        this.client = new WebSocketClient({
            autoReconnect: true,
            debug: true
        });

        // 设置事件处理器
        this.client.setEventHandlers({
            onMessage: (msg: string) => this.addMessage(msg),
            onStateChange: (state: ConnectionState) => this.updateStatus(state)
        });

        // 连接到服务器
        this.client.connect(url, token, deviceType);
    }

    public disconnect() {
        this.client?.close();
        this.client = null;
    }

    public sendMessage() {
        const messageInput = document.getElementById('messageInput') as HTMLInputElement;
        const message = messageInput.value;
        
        if (!message.trim()) {
            alert('请输入消息内容');
            return;
        }

        if (this.client?.send(message)) {
            console.log("消息发送成功");
            messageInput.value = '';
            this.addMessage(message, 'sent');
        } else {
            console.log("消息发送失败");
            alert('发送失败，请检查连接状态');
        }
    }
}

// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new WebSocketDemo();
}); 