// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserModel _$UserModelFromJson(Map<String, dynamic> json) => $checkedCreate(
      'UserModel',
      json,
      ($checkedConvert) {
        $checkKeys(
          json,
          allowedKeys: const [
            'id',
            'email',
            'role',
            'full_name',
            'enrollment_number',
            'created_at',
            'updated_at'
          ],
        );
        final val = UserModel(
          id: $checkedConvert('id', (v) => v as String),
          email: $checkedConvert('email', (v) => v as String),
          role:
              $checkedConvert('role', (v) => $enumDecode(_$UserRoleEnumMap, v)),
          fullName: $checkedConvert('full_name', (v) => v as String),
          enrollmentNumber:
              $checkedConvert('enrollment_number', (v) => v as String?),
          createdAt: $checkedConvert('created_at', (v) => (v as num).toInt()),
          updatedAt: $checkedConvert('updated_at', (v) => (v as num).toInt()),
        );
        return val;
      },
      fieldKeyMap: const {
        'fullName': 'full_name',
        'enrollmentNumber': 'enrollment_number',
        'createdAt': 'created_at',
        'updatedAt': 'updated_at'
      },
    );

Map<String, dynamic> _$UserModelToJson(UserModel instance) => <String, dynamic>{
      'id': instance.id,
      'email': instance.email,
      'role': _$UserRoleEnumMap[instance.role]!,
      'full_name': instance.fullName,
      if (instance.enrollmentNumber case final value?)
        'enrollment_number': value,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

const _$UserRoleEnumMap = {
  UserRole.student: 'student',
  UserRole.teacher: 'teacher',
};
