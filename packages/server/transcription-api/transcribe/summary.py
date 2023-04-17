import concurrent

from langchain.llms import Replicate
from langchain import PromptTemplate
from langchain.chains.summarize import load_summarize_chain
from langchain.docstore.document import Document

SUMMARISATION_AI_MODEL = "daanelson/flan-t5-xl:11d370d65d0040982f8435620af630b5965f7529d96494ab252b2ebb621e3169"


def summarise(transcript):
    print("inside summarise")
    dialogue = ""
    transcript_documents = []

    for line in transcript:
        dialogue +=  line['text'] + "\n"

        if len(dialogue.split(" ")) > 1000:
            transcript_documents.append(Document(page_content=dialogue))
            dialogue = ""

    if dialogue != "":
        transcript_documents.append(Document(page_content=dialogue))

    print("Transcript split into " + str(len(transcript_documents)) + " documents")
    prompt = PromptTemplate(
      template="Below is a transcribed video recording of a business meeting. Each new line indicates a change of speaker, the name of the speakers are not given. Please generate a summary of no more than 50 words of the most important points. Only use information contained below.\n\n{text}\n\nSummary:",
      input_variables=["text"],
    )

    merge = PromptTemplate(
      template="Below is a collection of summaries. Each line is a separate summary. Please combine them to make an concise overall summary of no more than 100 words that summarise the most important facts and action points. Only use information contained below. Filter out duplications.\n\n{text}\n\nOverall Summary:",
      input_variables=["text"],
    )

    output = {}
    try_count = 0
    while True:
        with concurrent.futures.ThreadPoolExecutor(max_workers=1) as executor:
            future = executor.submit(try_summary, merge, prompt, transcript_documents)
            try:
                output = future.result(timeout=(5*60)+(2.5*60*len(transcript_documents)))
                break
            except concurrent.futures.TimeoutError:
                print("Timeout Building Summary, retrying")
                try_count += 1
                if try_count > 3:
                    print("Timeout Building Summary, giving up")
                    break
                continue

    for step in output['intermediate_steps']:
        print(step)

    #output, steps = llm_chain.run(transcript_documents)
    result = output['output_text']
    print("result: " + result)
    return result


def try_summary(merge, prompt, transcript_documents):
    print("inside try_summary")
    llm = Replicate(model=SUMMARISATION_AI_MODEL,
                    input={"max_length": 200}, model_kwargs={"temperature": 0, "max_length": 200})
    llm_chain = load_summarize_chain(llm, map_prompt=prompt, combine_prompt=merge, chain_type="map_reduce",
                                     return_intermediate_steps=True)
    output = llm_chain({"input_documents": transcript_documents}, return_only_outputs=True)
    print("try_summary output: " + str(output))
    return output
