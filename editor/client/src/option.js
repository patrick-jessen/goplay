import {h} from "preact";

export default function Option({text, selected, options, onSelect, disabled}) {
  return (
    <tr>
      <td>{text}</td>
      <td>
        <select onChange={(e)=>onSelect(e.target.value)} disabled={disabled}>
          {
            options.map(o => (
              <option selected={o == selected}>{o}</option>
            ))
          }
        </select>
      </td>
    </tr>
  );
}