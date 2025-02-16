import 'dart:convert';
import 'package:json_annotation/json_annotation.dart';

part 'device_binding_model.g.dart';

@JsonSerializable()
class DeviceBindingModel {
  final String id;
  @JsonKey(name: 'user_id')
  final String userId;
  @JsonKey(name: 'device_id')
  final String deviceId;
  @JsonKey(name: 'device_info')
  final Map<String, dynamic> deviceInfo;
  @JsonKey(name: 'bound_at')
  final int boundAt;
  @JsonKey(name: 'created_at')
  final int createdAt;
  @JsonKey(name: 'updated_at')
  final int updatedAt;

  const DeviceBindingModel({
    required this.id,
    required this.userId,
    required this.deviceId,
    required this.deviceInfo,
    required this.boundAt,
    required this.createdAt,
    required this.updatedAt,
  });

  factory DeviceBindingModel.fromJson(Map<String, dynamic> json) =>
      _$DeviceBindingModelFromJson(json);

  Map<String, dynamic> toJson() => _$DeviceBindingModelToJson(this);

  Map<String, dynamic> toDatabase() => {
        'id': id,
        'user_id': userId,
        'device_id': deviceId,
        'device_info': jsonEncode(deviceInfo),
        'bound_at': boundAt,
        'created_at': createdAt,
        'updated_at': updatedAt,
      };

  factory DeviceBindingModel.fromDatabase(Map<String, dynamic> data) =>
      DeviceBindingModel(
        id: data['id'] as String,
        userId: data['user_id'] as String,
        deviceId: data['device_id'] as String,
        deviceInfo:
            jsonDecode(data['device_info'] as String) as Map<String, dynamic>,
        boundAt: data['bound_at'] as int,
        createdAt: data['created_at'] as int,
        updatedAt: data['updated_at'] as int,
      );

  DeviceBindingModel copyWith({
    String? id,
    String? userId,
    String? deviceId,
    Map<String, dynamic>? deviceInfo,
    int? boundAt,
    int? createdAt,
    int? updatedAt,
  }) {
    return DeviceBindingModel(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      deviceId: deviceId ?? this.deviceId,
      deviceInfo: deviceInfo ?? this.deviceInfo,
      boundAt: boundAt ?? this.boundAt,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
