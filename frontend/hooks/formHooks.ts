import { useCallback, useState, ChangeEvent } from "react";
import { Task } from "../src/generated/graphql";

function useTextInput(initialValue: string | undefined) {
  const [value, setValue] = useState(initialValue ?? "");

  const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setValue(event.target.value);
  }, []);
  return { inputProps: { value, onChange: handleChange }, setValue };
}

export function useTaskFields(initialTask?: Task) {
  const { inputProps: titleProps, setValue: setTitle } = useTextInput(
    initialTask?.title
  );
  const { inputProps: notesProps, setValue: setNotes } = useTextInput(
    initialTask?.note
  );

  const clearValue = useCallback(() => {
    setTitle("");
    setNotes("");
  }, [setNotes, setTitle]);

  return {
    titleProps,
    notesProps,
    clearValue
  };
}
