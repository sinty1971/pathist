'use client';

import React from 'react';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { TextField } from '@mui/material';
import { ja } from 'date-fns/locale';


interface CalendarPickerProps {
  value: string;
  onChange: (dateString: string) => void;
  placeholder: string;
  disabled?: boolean;
  minDate?: string;
}

export const CalendarPicker: React.FC<CalendarPickerProps> = ({
  value,
  onChange,
  placeholder,
  disabled = false,
  minDate
}) => {
  // value文字列をDateオブジェクトに変換
  const selectedDate = value ? new Date(value) : null;
  const [internalDate, setInternalDate] = React.useState<Date | null>(selectedDate);

  // valueが外部から変更された時に内部状態を同期
  React.useEffect(() => {
    setInternalDate(selectedDate);
  }, [value]);

  const handleDateChange = (date: Date | null) => {
    // 内部状態のみ更新（表示用）
    setInternalDate(date);
    // 親コンポーネントには通知しない
  };

  const handleDateAccept = (date: Date | null) => {
    // onAcceptは日付が確定した時のみ呼ばれる
    if (date) {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      onChange(`${year}-${month}-${day}`);
    } else {
      onChange('');
    }
  };

  // minDate文字列をDateオブジェクトに変換
  const minDateObj = minDate ? new Date(minDate) : undefined;

  return (
    <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ja}>
      <DatePicker
        value={internalDate}
        onChange={handleDateChange}
        onAccept={handleDateAccept}
        disabled={disabled}
        minDate={minDateObj}
        format="yyyy年MM月dd日"
        enableAccessibleFieldDOMStructure={false}
        views={['year', 'month', 'day']}
        openTo="day"
        reduceAnimations={true}
        closeOnSelect={false}
        slots={{
          textField: TextField,
        }}
        slotProps={{
          textField: {
            placeholder: placeholder,
            size: 'small',
            fullWidth: true,
            sx: {
              '& .MuiOutlinedInput-root': {
                fontSize: '0.875rem',
                borderRadius: '6px',
                '&:hover .MuiOutlinedInput-notchedOutline': {
                  borderColor: '#3b82f6',
                },
                '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                  borderColor: '#3b82f6',
                  boxShadow: '0 0 0 3px rgba(59, 130, 246, 0.1)',
                },
              },
              '& .MuiInputBase-input': {
                padding: '8px 12px',
              },
            },
          },
        }}
      />
    </LocalizationProvider>
  );
};