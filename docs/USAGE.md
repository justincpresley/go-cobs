# Usage

## Configuration

The following are variables of the `Config` structure. All variables are applicable to all types. If you are frightened by the options, feel free to use the suggested defaults or create an empty structure (which acts like a very simple implementation).

* [X] **SpecialByte**:
  * Description: The particular byte value to be "encoded away".
  * Possible Values: Any hex value such as `0x00`.
  * Suggested Default: No default.
* [X] **Delimiter**:
  * Description: Whether to include a delimiter: the special byte value placed at the end of the encoded slice to mark the end.
  * Possible Values: The two boolean values: `true` and `false`.
  * Suggested Default: No default.
* [X] **Type**:
  * Description: The Type of COBS that is used for all API functions.
  * Possible Values: Any `Type` value.
  * Suggested Default: `Native`.
* [X] **EndingSave**:
  * Description: Whether to save a byte when encountering "max" flag(s) occurring at the end of the slice.
  * Possible Values: The two boolean values: `true` and `false`.
  * Suggested Default: `true`.
* [ ] **Reverse**:
  * Description: Place the flag at the end of the chunk rather than before effectively reversing the process - `RCOBS`. This allows encoding with no lookahead (making it easier to encode) but enforces decoding to be done in reverse killing the ability to stream decode.
  * Possible Values: The two boolean values: `true` and `false`.
  * Suggested Default: `false`.

## Types

The following are the types (extensions) that are/will be included. You can think of types as different versions COBS. Each has their own advantages, disadvantages, and use case.

Individual Types:

* [X] **Native** ``(COBS)``:
  * Description: This is the natural and default algorithm that was first presented.
  * Pros: Relatively stable overhead. The easiest to implement and possibly the most performant.
  * Cons: Nothing is bad about the default! But it could be improved...
* [X] **Reduced** ``(COBS/R)``:
  * Description: Saves a byte by replacing the last flag with the last character if is appropriate!
  * Pros: Potentially saves a byte of overhead. Encoding possibly generates no overhead.
  * Cons: A massive reduction in coverage from flag-based verification.
* [X] **PairElimination** ``(COBS/PE)``:
  * Description: Incorporate flags to represent a "pair" of special bytes.
  * Pros: A common reduction in overhead. Best case can be almost half of the size. Good for embedded systems.
  * Cons: An increase in theoretical worse case (maximum overhead) than `Native`.
* [ ] **RunElimination** ``(COBS/RE)``:
  * Description: Incorporate flags to represent a "run" of special bytes.
  * Pros: A rare but massive reduction in overhead. Good for embedded systems.
  * Cons: An increase in theoretical worse case (maximum overhead) than `Native`.

Combined Types:

* [ ] **PairAndRun** ``(COBS/PAR)``:
  * Description: Combined `PairElimination` and `RunElimination`.
  * Pros: Objectively achieving an optimal balance. Pros of both types.
  * Cons: Cons of both types.