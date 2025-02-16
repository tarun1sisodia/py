import 'package:smart_campus/domain/entities/course.dart';

abstract class CourseRepository {
  Future<List<Course>> getCourses({
    String? department,
    int? yearOfStudy,
    int? semester,
  });

  Future<Course> getCourseById(String id);

  Future<Course> getCourseByCode(String courseCode);

  Future<List<Course>> getTeacherCourses(String teacherId);

  Future<List<Course>> getStudentCourses(String studentId);

  Future<Course> createCourse(Course course);

  Future<Course> updateCourse(Course course);

  Future<void> deleteCourse(String id);

  Future<void> assignTeacherToCourse({
    required String teacherId,
    required String courseId,
    required String academicYear,
  });

  Future<void> removeTeacherFromCourse({
    required String teacherId,
    required String courseId,
    required String academicYear,
  });

  Future<List<Course>> searchCourses(String query);

  Future<bool> isCourseCodeUnique(String courseCode);
}
