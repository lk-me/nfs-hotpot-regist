#include <stdio.h>
#include <stdlib.h>
#include <dbus/dbus.h>

#define DBUS_DEST   "org.nfs.HotpotRegist1"
#define DBUS_PATH   "/org/nfs/HotpotRegist1"
#define DBUS_IFC    "org.nfs.HotpotRegist1"

int setBooleanValue(DBusConnection *connection) {

    fprintf(stderr, "使用方法：真假 [1|0]\n");
    int boolValue = 0;
    scanf("%d", &boolValue);

    DBusError error;
    DBusMessage *message, *reply;

    // 初始化DBus错误对象
    dbus_error_init(&error);

    // 创建一个设置布尔值的DBus消息
    message = dbus_message_new_method_call(
        DBUS_DEST,  // 目标接口的名称
        DBUS_PATH,    // 目标对象的路径
        DBUS_IFC,  // 目标接口的名称
        "SaveRegistInfo"         // 目标方法的名称
    );
    if (message == NULL) {
        fprintf(stderr, "DBus消息创建失败\n");
        return 1;
    }

    // 添加布尔参数到DBus消息中
    if (!dbus_message_append_args(message, DBUS_TYPE_BOOLEAN, &boolValue, DBUS_TYPE_INVALID)) {
        fprintf(stderr, "布尔参数添加失败\n");
        dbus_message_unref(message);
        return 1;
    }

    // 发送DBus消息并等待回复
    reply = dbus_connection_send_with_reply_and_block(connection, message, -1, &error);
    if (dbus_error_is_set(&error)) {
        fprintf(stderr, "DBus消息发送错误: %s\n", error.message);
        dbus_error_free(&error);
        dbus_message_unref(message);
        return 1;
    }

    // 释放DBus消息
    dbus_message_unref(message);

    // 处理返回的DBus回复
    if (dbus_message_get_type(reply) == DBUS_MESSAGE_TYPE_ERROR) {
        fprintf(stderr, "DBus方法调用错误: %s\n", dbus_message_get_error_name(reply));
        dbus_message_unref(reply);
        return 1;
    }

    // 释放DBus回复
    dbus_message_unref(reply);

    printf("布尔值已成功设置为 %d\n", boolValue);

    return 0;
}

// 获取布尔值
int getBooleanValue(DBusConnection *connection) {
    DBusError error;
    DBusMessage *message, *reply;
    dbus_bool_t returnValue;

    // 初始化DBus错误对象
    dbus_error_init(&error);

    // 创建一个获取布尔值的DBus消息
    message = dbus_message_new_method_call(
        DBUS_DEST,  // 目标接口的名称
        DBUS_PATH,    // 目标对象的路径
        DBUS_IFC,  // 目标接口的名称
        "ReadRegistInfo"         // 目标方法的名称
    );
    if (message == NULL) {
        fprintf(stderr, "DBus消息创建失败\n");
        return 1;
    }

    // 发送DBus消息并等待回复
    reply = dbus_connection_send_with_reply_and_block(connection, message, -1, &error);
    if (dbus_error_is_set(&error)) {
        fprintf(stderr, "DBus消息发送错误: %s\n", error.message);
        dbus_error_free(&error);
        dbus_message_unref(message);
        return 1;
    }

    // 释放DBus消息
    dbus_message_unref(message);

    // 处理返回的DBus回复
    if (dbus_message_get_type(reply) == DBUS_MESSAGE_TYPE_ERROR) {
        fprintf(stderr, "DBus方法调用错误: %s\n", dbus_message_get_error_name(reply));
        dbus_message_unref(reply);
        return 1;
    }

    // 从DBus回复中解析布尔返回值
    if (!dbus_message_get_args(reply, &error, DBUS_TYPE_BOOLEAN, &returnValue, DBUS_TYPE_INVALID)) {
        fprintf(stderr, "布尔返回值解析失败: %s\n", error.message);
        dbus_error_free(&error);
        dbus_message_unref(reply);
        return 1;
    }

    // 释放DBus回复
    dbus_message_unref(reply);

    printf("布尔返回值: %s\n",returnValue ? "True" : "False");

    return 0;
}

int main(int argc, char *argv[]) {
    if (argc != 2) {
        fprintf(stderr, "使用方法：程序名 [1 设置|2 读取]\n");
        return 1;
    }

    int option = atoi(argv[1]);

    DBusError error;
    DBusConnection *connection;
    // 初始化DBus错误对象
    dbus_error_init(&error);
    // 连接到DBus系统总线
    connection = dbus_bus_get(DBUS_BUS_SYSTEM, &error);
    if (option == 1) {
        if (0 != setBooleanValue(connection)){
             return 1;
        }
    } else {
        if (0 != getBooleanValue(connection)) {
             return 1;
        }
    }
    dbus_connection_unref(connection);
}