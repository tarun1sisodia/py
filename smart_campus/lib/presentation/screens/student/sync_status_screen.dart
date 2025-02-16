import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:smart_campus/config/theme.dart';
import 'package:smart_campus/core/services/sync_service.dart';
import 'package:smart_campus/domain/entities/attendance_record.dart';

class SyncStatusScreen extends StatefulWidget {
  const SyncStatusScreen({super.key});

  @override
  State<SyncStatusScreen> createState() => _SyncStatusScreenState();
}

class _SyncStatusScreenState extends State<SyncStatusScreen> {
  List<AttendanceRecord> _pendingRecords = [];
  List<AttendanceRecord> _failedRecords = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadRecords();
  }

  Future<void> _loadRecords() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final pendingRecords =
          await context.read<SyncService>().getPendingRecords();
      final failedRecords =
          await context.read<SyncService>().getFailedRecords();

      setState(() {
        _pendingRecords = pendingRecords;
        _failedRecords = failedRecords;
        _isLoading = false;
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Failed to load records: ${e.toString()}'),
          backgroundColor: Colors.red,
        ),
      );
      setState(() {
        _isLoading = false;
      });
    }
  }

  Future<void> _retryFailedRecord(AttendanceRecord record) async {
    try {
      await context
          .read<SyncService>()
          .retryFailedRecord(record.sessionId, record.studentId);
      _loadRecords();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Failed to retry record: ${e.toString()}'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Sync Status'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _loadRecords,
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _loadRecords,
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildSyncSummary(),
                    const SizedBox(height: 24),
                    _buildPendingRecords(),
                    const SizedBox(height: 24),
                    _buildFailedRecords(),
                  ],
                ),
              ),
            ),
    );
  }

  Widget _buildSyncSummary() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Sync Summary',
              style: AppTheme.titleLarge,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _buildSummaryItem(
                    'Pending',
                    _pendingRecords.length.toString(),
                    Icons.pending,
                    Colors.orange,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: _buildSummaryItem(
                    'Failed',
                    _failedRecords.length.toString(),
                    Icons.error,
                    Colors.red,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSummaryItem(
    String label,
    String value,
    IconData icon,
    Color color,
  ) {
    return Container(
      padding: const EdgeInsets.all(16.0),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color),
      ),
      child: Column(
        children: [
          Icon(icon, color: color),
          const SizedBox(height: 8),
          Text(
            value,
            style: AppTheme.headlineMedium.copyWith(color: color),
          ),
          Text(
            label,
            style: AppTheme.bodyMedium.copyWith(color: color),
          ),
        ],
      ),
    );
  }

  Widget _buildPendingRecords() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Pending Records',
          style: AppTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        if (_pendingRecords.isEmpty)
          const Center(
            child: Text('No pending records'),
          )
        else
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _pendingRecords.length,
            itemBuilder: (context, index) {
              final record = _pendingRecords[index];
              return _buildRecordCard(record);
            },
          ),
      ],
    );
  }

  Widget _buildFailedRecords() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Failed Records',
          style: AppTheme.titleLarge,
        ),
        const SizedBox(height: 16),
        if (_failedRecords.isEmpty)
          const Center(
            child: Text('No failed records'),
          )
        else
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _failedRecords.length,
            itemBuilder: (context, index) {
              final record = _failedRecords[index];
              return _buildRecordCard(record, showRetry: true);
            },
          ),
      ],
    );
  }

  Widget _buildRecordCard(AttendanceRecord record, {bool showRetry = false}) {
    return Card(
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: record.status == VerificationStatus.pending
              ? Colors.orange
              : Colors.red,
          child: Icon(
            record.status == VerificationStatus.pending
                ? Icons.pending
                : Icons.error,
            color: Theme.of(context).colorScheme.onPrimary,
          ),
        ),
        title: Text(
          'Session ID: ${record.sessionId}',
          style: AppTheme.bodyLarge,
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Marked at: ${record.markedAt.toString()}',
              style: AppTheme.bodyMedium,
            ),
            if (record.rejectionReason != null)
              Text(
                'Error: ${record.rejectionReason}',
                style: AppTheme.bodyMedium.copyWith(color: Colors.red),
              ),
          ],
        ),
        trailing: showRetry
            ? IconButton(
                icon: const Icon(Icons.refresh),
                onPressed: () => _retryFailedRecord(record),
              )
            : null,
      ),
    );
  }
}
