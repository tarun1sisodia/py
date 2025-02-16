// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'course_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CourseModel _$CourseModelFromJson(Map<String, dynamic> json) => $checkedCreate(
      'CourseModel',
      json,
      ($checkedConvert) {
        $checkKeys(
          json,
          allowedKeys: const [
            'id',
            'code',
            'name',
            'department',
            'year_of_study',
            'semester',
            'created_at',
            'updated_at'
          ],
        );
        final val = CourseModel(
          id: $checkedConvert('id', (v) => v as String),
          code: $checkedConvert('code', (v) => v as String),
          name: $checkedConvert('name', (v) => v as String),
          department: $checkedConvert('department', (v) => v as String),
          yearOfStudy:
              $checkedConvert('year_of_study', (v) => (v as num).toInt()),
          semester: $checkedConvert('semester', (v) => (v as num).toInt()),
          createdAt: $checkedConvert('created_at', (v) => (v as num).toInt()),
          updatedAt: $checkedConvert('updated_at', (v) => (v as num).toInt()),
        );
        return val;
      },
      fieldKeyMap: const {
        'yearOfStudy': 'year_of_study',
        'createdAt': 'created_at',
        'updatedAt': 'updated_at'
      },
    );

Map<String, dynamic> _$CourseModelToJson(CourseModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'code': instance.code,
      'name': instance.name,
      'department': instance.department,
      'year_of_study': instance.yearOfStudy,
      'semester': instance.semester,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };
