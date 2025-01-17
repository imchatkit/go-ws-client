即时通讯 SDK 命名规范
1. 包名/模块名规范
各语言包名应遵循以下规范：
iOS/Swift: IMKit
Android/Java: com.imai.im
Flutter: im_kit
Go: github.com/imai/im
C++: imai::im
TypeScript: @imai/im-sdk
2. 核心类/接口命名
2.1 客户端实例
统一使用 IMClient 作为客户端实例名称：
// 各平台示例
iOS/Swift: IMClient.shared
Android: IMClient.getInstance()
Flutter: IMClient.instance
C++: IMClient::getInstance()
TypeScript: new IMClient()
2.2 连接管理器
IMConnectionManager - 负责 WebSocket 连接的建立、维护、重连等：
interface IMConnectionManager {
// 基础连接方法
connect(token: string): Promise<void>;
disconnect(): void;
reconnect(): void;
// 状态查询
getConnectionState(): IMConnectionState;
isConnected(): boolean;
// 事件回调
onConnectionStateChanged(state: IMConnectionState): void;
onConnected(): void;
onDisconnected(reason: IMDisconnectReason): void;
onReconnecting(): void;
onReconnected(): void;
}
2.3 消息管理器
IMMessageManager - 负责消息的发送、接收、存储等：
interface IMMessageManager {
// 发送消息
sendMessage(message: IMMessage): Promise<void>;
sendTextMessage(text: string): Promise<void>;
sendImageMessage(image: IMImage): Promise<void>;
// 消息操作
recallMessage(messageId: string): Promise<void>;
resendMessage(messageId: string): Promise<void>;
// 消息回调
onMessageReceived(message: IMMessage): void;
onMessageSent(messageId: string): void;
onMessageDelivered(messageId: string): void;
onMessageRead(messageId: string): void;
onMessageRecalled(messageId: string): void;
}
2.4 会话管理器
IMConversationManager - 负责会话的管理：
interface IMConversationManager {
// 会话操作
getConversation(conversationId: string): IMConversation;
getConversationList(): Promise<IMConversation[]>;
deleteConversation(conversationId: string): Promise<void>;
// 未读消息
getUnreadCount(): number;
markAsRead(conversationId: string): void;
// 事件回调
onConversationUpdated(conversation: IMConversation): void;
onUnreadCountChanged(count: number): void;
}
3. 数据模型命名
3.1 消息模型
interface IMMessage {
messageId: string;
conversationId: string;
senderId: string;
content: IMMessageContent;
timestamp: number;
status: IMMessageStatus;
}
interface IMMessageContent {
type: IMMessageType;
text?: string;
image?: IMImage;
file?: IMFile;
custom?: any;
}
3.2 会话模型
interface IMConversation {
conversationId: string;
type: IMConversationType;
unreadCount: number;
lastMessage: IMMessage;
updatedAt: number;
}
4. 枚举类型命名
4.1 连接状态
enum IMConnectionState {
DISCONNECTED = "DISCONNECTED",
CONNECTING = "CONNECTING",
CONNECTED = "CONNECTED",
RECONNECTING = "RECONNECTING"
}
4.2 消息状态
enum IMMessageStatus {
SENDING = "SENDING",
SENT = "SENT",
DELIVERED = "DELIVERED",
READ = "READ",
FAILED = "FAILED"
}
4.3 消息类型
enum IMMessageType {
TEXT = "TEXT",
IMAGE = "IMAGE",
FILE = "FILE",
CUSTOM = "CUSTOM"
}
5. 错误处理命名
5.1 错误码
enum IMErrorCode {
// 连接相关 1xxxx
CONNECTION_TIMEOUT = 10001,
TOKEN_EXPIRED = 10002,
NETWORK_UNAVAILABLE = 10003,
// 消息相关 2xxxx
MESSAGE_TOO_LARGE = 20001,
MESSAGE_SEND_FAILED = 20002,
// 会话相关 3xxxx
CONVERSATION_NOT_FOUND = 30001
}
5.2 错误回调
interface IMErrorCallback {
onError(error: IMError): void;
onConnectionError(error: IMConnectionError): void;
onMessageError(error: IMMessageError): void;
}
6. 配置选项命名
interface IMClientConfig {
appId: string;
deviceId: string;
debugMode: boolean;
heartbeatInterval: number;
reconnectStrategy: IMReconnectStrategy;
storage: IMStorageType;
}
7. 使用示例
// 初始化客户端
const client = new IMClient({
appId: "your-app-id",
deviceId: "device-id",
debugMode: true
});
// 设置回调
client.onConnectionStateChanged = (state) => {
console.log("Connection state:", state);
};
client.onMessageReceived = (message) => {
console.log("New message:", message);
};
// 连接服务器
await client.connect("your-token");
// 发送消息
await client.messageManager.sendTextMessage("Hello, World!");
8. 跨平台支持规范

### 8.1 平台适配器

```typescript
// 定义平台适配器接口
interface IMPlatformAdapter {
    // 网络相关
    createWebSocket(url: string): IMWebSocket;
    isNetworkAvailable(): boolean;

    // 存储相关
    getStorage(): IMStorage;
    
    // 设备信息
    getDeviceInfo(): IMDeviceInfo;
    
    // 通知
    showNotification(message: IMMessage): void;
}

// 各平台实现
class WebAdapter implements IMPlatformAdapter { ... }
class ReactNativeAdapter implements IMPlatformAdapter { ... }
class ElectronAdapter implements IMPlatformAdapter { ... }
```

### 8.2 存储接口

```typescript
interface IMStorage {
    // 通用存储接口
    get(key: string): Promise<any>;
    set(key: string, value: any): Promise<void>;
    remove(key: string): Promise<void>;
    clear(): Promise<void>;
}

// 平台特定实现
class WebStorage implements IMStorage { ... }      // localStorage
class RNStorage implements IMStorage { ... }       // AsyncStorage
class ElectronStorage implements IMStorage { ... } // electron-store
```

### 8.3 初始化配置

```typescript
interface IMClientConfig {
    // ... 原有配置 ...
    
    // 平台相关配置
    platform: IMPlatformType;
    adapter?: IMPlatformAdapter;    // 自定义适配器
    storage?: IMStorage;            // 自定义存储
    
    // 平台特定选项
    web?: {
        enableNotification: boolean;
    };
    reactNative?: {
        enableBackgroundMessage: boolean;
    };
    electron?: {
        tray: boolean;
        windowNotification: boolean;
    };
}
```

### 8.4 使用示例

```typescript
// Web 环境 (Vue/React)
const client = new IMClient({
    platform: IMPlatformType.WEB,
    web: {
        enableNotification: true
    }
});

// React Native 环境
const client = new IMClient({
    platform: IMPlatformType.REACT_NATIVE,
    reactNative: {
        enableBackgroundMessage: true
    }
});

// Electron 环境
const client = new IMClient({
    platform: IMPlatformType.ELECTRON,
    electron: {
        tray: true,
        windowNotification: true
    }
});
```

### 8.5 目录结构

```
@imai/im-sdk/
├── src/
│   ├── core/                 # 核心功能
│   ├── platforms/           # 平台适配
│   │   ├── web/
│   │   ├── react-native/
│   │   └── electron/
│   └── utils/
├── dist/
│   ├── web/
│   ├── react-native/
│   └── electron/
└── examples/
     ├── vue/
     ├── react/
     ├── react-native/
     └── electron/
```