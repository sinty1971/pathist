import { useState, useEffect } from 'react';
// import { api } from '../api/client';

interface DateEditModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  projectName: string;
  currentStartDate: string;
  currentEndDate: string;
  onSuccess: (projectId: string, startDate: string, endDate: string) => void;
}

const DateEditModal = ({ 
  isOpen, 
  onClose, 
  projectId, 
  projectName, 
  currentStartDate, 
  currentEndDate, 
  onSuccess 
}: DateEditModalProps) => {
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [isUpdating, setIsUpdating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      // Convert ISO strings to local date strings for input
      const formatDateForInput = (dateString: string) => {
        if (!dateString) return '';
        try {
          const date = new Date(dateString);
          // Use local time zone methods to avoid UTC conversion issues
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, '0');
          const day = String(date.getDate()).padStart(2, '0');
          return `${year}-${month}-${day}`;
        } catch {
          return '';
        }
      };
      
      setStartDate(formatDateForInput(currentStartDate));
      setEndDate(formatDateForInput(currentEndDate));
      setError(null);
    }
  }, [isOpen, currentStartDate, currentEndDate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!startDate || !endDate) {
      setError('開始日と終了日の両方を入力してください');
      return;
    }

    if (new Date(startDate) > new Date(endDate)) {
      setError('開始日は終了日より前である必要があります');
      return;
    }

    setIsUpdating(true);
    setError(null);

    try {
      // Send dates as simple YYYY-MM-DD format to avoid timezone conversion issues
      // The backend will handle these as local dates
      const startDateStr = `${startDate}T00:00:00`;
      const endDateStr = `${endDate}T23:59:59`;
      
      // 編集された日付を親コンポーネントに通知
      onSuccess(projectId, startDateStr, endDateStr);
      onClose();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '日付の更新に失敗しました';
      setError(errorMessage);
    } finally {
      setIsUpdating(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>プロジェクト日付編集</h2>
          <button className="close-button" onClick={onClose}>
            ×
          </button>
        </div>
        
        <div className="modal-body">
          <p><strong>プロジェクト:</strong> {projectName}</p>
          <p><strong>プロジェクトID:</strong> {projectId}</p>
          
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="start-date">開始日:</label>
              <input
                id="start-date"
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                required
                className="date-input"
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="end-date">終了日:</label>
              <input
                id="end-date"
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                required
                className="date-input"
              />
            </div>
            
            {error && <div className="error-message">{error}</div>}
            
            <div className="form-actions">
              <button 
                type="button" 
                onClick={onClose}
                className="cancel-button"
                disabled={isUpdating}
              >
                キャンセル
              </button>
              <button 
                type="submit" 
                className="submit-button"
                disabled={isUpdating}
              >
                {isUpdating ? '更新中...' : '更新'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default DateEditModal;