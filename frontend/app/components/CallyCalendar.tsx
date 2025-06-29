import { useState } from 'react';

interface CallyCalendarProps {
  className?: string;
  onDateSelect?: (date: string) => void;
  selectedDate?: string; // 現在選択されている日付 (YYYY-MM-DD形式)
}

export default function CallyCalendar({ className = '', onDateSelect, selectedDate }: CallyCalendarProps) {
  // 選択されている日付があればその月を表示、なければ今日の月を表示
  const getInitialDate = () => {
    if (selectedDate) {
      const parts = selectedDate.split('-');
      if (parts.length === 3) {
        return new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]));
      }
    }
    return new Date();
  };
  
  const [currentDate, setCurrentDate] = useState(getInitialDate());
  const [showYearPicker, setShowYearPicker] = useState(false);
  const [showMonthPicker, setShowMonthPicker] = useState(false);

  const handleMonthSelect = (month: number) => {
    const newDate = new Date(currentDate);
    newDate.setMonth(month);
    setCurrentDate(newDate);
    setShowMonthPicker(false);
  };

  const handleDateClick = (day: number) => {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();
    
    // ローカルタイムゾーンでの日付文字列を作成（UTCを使わない）
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
    
    if (onDateSelect) {
      onDateSelect(dateStr);
    }
  };

  const handleYearSelect = (year: number) => {
    const newDate = new Date(currentDate);
    newDate.setFullYear(year);
    setCurrentDate(newDate);
    setShowYearPicker(false);
  };

  const renderYearPicker = () => {
    const currentYear = currentDate.getFullYear();
    const years = [];
    for (let i = currentYear - 3; i <= currentYear + 3; i++) {
      years.push(
        <button
          key={i}
          className={`year-option ${i === currentYear ? 'selected' : ''}`}
          onClick={() => handleYearSelect(i)}
          type="button"
        >
          {i}
        </button>
      );
    }
    return years;
  };

  const renderMonthPicker = () => {
    const months = [
      '1月', '2月', '3月', '4月', '5月', '6月',
      '7月', '8月', '9月', '10月', '11月', '12月'
    ];
    const currentMonth = currentDate.getMonth();
    
    return months.map((monthName, index) => (
      <button
        key={index}
        className={`month-option ${index === currentMonth ? 'selected' : ''}`}
        onClick={() => handleMonthSelect(index)}
        type="button"
      >
        {monthName}
      </button>
    ));
  };

  const renderCalendar = () => {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const firstDayOfWeek = firstDay.getDay();
    const daysInMonth = lastDay.getDate();

    const days = [];
    const today = new Date();
    
    // 選択されている日付を解析
    let selectedDay = 0;
    if (selectedDate) {
      const parts = selectedDate.split('-');
      if (parts.length === 3) {
        const selectedYear = parseInt(parts[0]);
        const selectedMonth = parseInt(parts[1]) - 1;
        const selectedDayNum = parseInt(parts[2]);
        if (selectedYear === year && selectedMonth === month) {
          selectedDay = selectedDayNum;
        }
      }
    }
    
    // 空白セルを追加
    for (let i = 0; i < firstDayOfWeek; i++) {
      days.push(<div key={`empty-${i}`} className="calendar-day empty"></div>);
    }

    // 日付セルを追加
    for (let day = 1; day <= daysInMonth; day++) {
      const isToday = today.getFullYear() === year && 
                     today.getMonth() === month && 
                     today.getDate() === day;
      const isSelected = day === selectedDay;
      
      let className = 'calendar-day';
      if (isSelected) {
        className += ' selected';
      } else if (isToday) {
        className += ' today';
      }
      
      days.push(
        <button
          key={day}
          className={className}
          onClick={() => handleDateClick(day)}
          type="button"
        >
          {day}
        </button>
      );
    }

    return days;
  };

  return (
    <div className={`simple-calendar ${className}`}>
      <style>{`
        .simple-calendar {
          width: 200px;
          max-width: 100%;
          border: 1px solid #e0e0e0;
          border-radius: 8px;
          background: white;
          padding: 8px;
          position: relative;
        }
        .calendar-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 8px;
        }
        .nav-button {
          background: none;
          border: none;
          cursor: pointer;
          padding: 2px 6px;
          border-radius: 4px;
          font-size: 12px;
        }
        .nav-button:hover {
          background-color: #f0f0f0;
        }
        .month-year {
          font-weight: 500;
          font-size: 12px;
        }
        .year-clickable {
          background: none;
          border: none;
          cursor: pointer;
          padding: 2px 4px;
          border-radius: 4px;
          text-decoration: underline;
        }
        .year-clickable:hover {
          background-color: #f0f0f0;
        }
        .year-picker {
          position: absolute;
          top: 40px;
          left: 50%;
          transform: translateX(-50%);
          background: white;
          border: 1px solid #ccc;
          border-radius: 8px;
          padding: 8px;
          box-shadow: 0 4px 12px rgba(0,0,0,0.15);
          z-index: 1000;
          display: flex;
          flex-direction: column;
          gap: 4px;
        }
        .year-option {
          background: none;
          border: none;
          cursor: pointer;
          padding: 6px 12px;
          border-radius: 4px;
          font-size: 12px;
          text-align: center;
          min-width: 60px;
        }
        .year-option:hover {
          background-color: #e3f2fd;
        }
        .year-option.selected {
          background-color: #1976d2;
          color: white;
        }
        .month-picker {
          position: absolute;
          top: 40px;
          left: 50%;
          transform: translateX(-50%);
          background: white;
          border: 1px solid #ccc;
          border-radius: 8px;
          padding: 8px;
          box-shadow: 0 4px 12px rgba(0,0,0,0.15);
          z-index: 1000;
          display: flex;
          flex-direction: column;
          gap: 4px;
          max-height: 180px;
          overflow-y: auto;
        }
        .month-option {
          background: none;
          border: none;
          cursor: pointer;
          padding: 6px 12px;
          border-radius: 4px;
          font-size: 12px;
          text-align: center;
          min-width: 60px;
        }
        .month-option:hover {
          background-color: #e3f2fd;
        }
        .month-option.selected {
          background-color: #1976d2;
          color: white;
        }
        .calendar-grid {
          display: grid;
          grid-template-columns: repeat(7, 1fr);
          gap: 2px;
        }
        .day-header {
          text-align: center;
          font-size: 11px;
          font-weight: 500;
          padding: 4px;
          color: #666;
        }
        .calendar-day {
          width: 24px;
          height: 24px;
          border: none;
          background: none;
          cursor: pointer;
          border-radius: 4px;
          font-size: 10px;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .calendar-day:hover {
          background-color: #e3f2fd;
        }
        .calendar-day.today {
          background-color: #e3f2fd;
          color: #1976d2;
          border: 1px solid #1976d2;
        }
        .calendar-day.selected {
          background-color: #1976d2;
          color: white;
        }
        .calendar-day.empty {
          cursor: default;
        }
      `}</style>
      
      <div className="calendar-header">
        <div className="month-year">
          <button 
            className="year-clickable"
            onClick={() => setShowYearPicker(!showYearPicker)}
            type="button"
          >
            {currentDate.getFullYear()}年
          </button>
          <button 
            className="year-clickable"
            onClick={() => setShowMonthPicker(!showMonthPicker)}
            type="button"
          >
            {currentDate.getMonth() + 1}月
          </button>
        </div>
      </div>

      {showYearPicker && (
        <div className="year-picker">
          {renderYearPicker()}
        </div>
      )}

      {showMonthPicker && (
        <div className="month-picker">
          {renderMonthPicker()}
        </div>
      )}

      <div className="calendar-grid">
        <div className="day-header">日</div>
        <div className="day-header">月</div>
        <div className="day-header">火</div>
        <div className="day-header">水</div>
        <div className="day-header">木</div>
        <div className="day-header">金</div>
        <div className="day-header">土</div>
        {renderCalendar()}
      </div>
    </div>
  );
}