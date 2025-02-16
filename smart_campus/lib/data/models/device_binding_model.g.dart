// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'device_binding_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

DeviceBindingModel _$DeviceBindingModelFromJson(Map<String, dynamic> json) =>
    $checkedCreate(
      'DeviceBindingModel',
      json,
      ($checkedConvert) {
        $checkKeys(
          json,
          allowedKeys: const [
            'id',
            'user_id',
            'device_id',
            'device_info',
            'bound_at',
            'created_at',
            'updated_at'
          ],
        );
        final val = DeviceBindingModel(
          id: $checkedConvert('id', (v) => v as String),
          userId: $checkedConvert('user_id', (v) => v as String),
          deviceId: $checkedConvert('device_id', (v) => v as String),
          deviceInfo:
              $checkedConvert('device_info', (v) => v as Map<String, dynamic>),
          boundAt: $checkedConvert('bound_at', (v) => (v as num).toInt()),
          createdAt: $checkedConvert('created_at', (v) => (v as num).toInt()),
          updatedAt: $checkedConvert('updated_at', (v) => (v as num).toInt()),
        );
        return val;
      },
      fieldKeyMap: const {
        'userId': 'user_id',
        'deviceId': 'device_id',
        'deviceInfo': 'device_info',
        'boundAt': 'bound_at',
        'createdAt': 'created_at',
        'updatedAt': 'updated_at'
      },
    );

Map<String, dynamic> _$DeviceBindingModelToJson(DeviceBindingModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'user_id': instance.userId,
      'device_id': instance.deviceId,
      'device_info': instance.deviceInfo,
      'bound_at': instance.boundAt,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };
