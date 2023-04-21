import concurrent
import json

from langchain.llms import HuggingFaceHub
from langchain.text_splitter import CharacterTextSplitter, RecursiveCharacterTextSplitter
from transformers import AutoTokenizer, pipeline

from langchain import PromptTemplate, LLMChain
from langchain.docstore.document import Document

SUMMARISATION_AI_MODEL = "knkarthick/MEETING_SUMMARY"
TOKEN_LIMIT = 800

def summarise(transcript):
    print(f"inside summarise:\n{json.dumps(transcript)}")
    tokenizer = AutoTokenizer.from_pretrained(SUMMARISATION_AI_MODEL)
    transcript_documents = split_text_into_documents(tokenizer, [x['text'] for x in transcript])

    print("Transcript split into " + str(len(transcript_documents)) + " documents")
    prompt = PromptTemplate(
      template="Below is a transcribed video recording of a business meeting. Each new line indicates a change of speaker, the name of the speakers are not given. Please generate a 50 word summary containing the most important points. Only use information contained below.\n\n{text}\n\nSummary:",
      input_variables=["text"],
    )

    merge = PromptTemplate(
      template="Below is a collection of summaries. Each line is a separate summary. Please combine all the summaries to make an concise overall summary of no more than 100 words. Only use information contained below. Filter out duplications.\n\n{text}\n\nOverall Summary:",
      input_variables=["text"],
    )

    output = {}
    summarizer = pipeline("summarization", model=SUMMARISATION_AI_MODEL, device_map="auto")

    run_count = 0
    while len(transcript_documents) > 1:
        transcript_summaries = summarizer([x.page_content for x in transcript_documents], batch_size=8, max_length=100, min_length=10, do_sample=False)
        for transcript_summary in transcript_summaries:
            print (f"Intermediate output run {run_count}: {json.dumps(transcript_summary)}")
        run_count += 1
        transcript_documents = split_text_into_documents(tokenizer, [x['summary_text'] for x in transcript_summaries])

    final_transcript = summarizer(transcript_documents[0].page_content, max_length=100, min_length=10, do_sample=False)
    output = {'output_text':final_transcript[0]['summary_text'], 'intermediate_steps': []}


    #output, steps = llm_chain.run(transcript_documents)
    result = output['output_text']
    print("result: " + result)
    return result


def split_text_into_documents(tokenizer, transcript):
    dialogue = ""
    transcript_documents = []
    text_splitter = RecursiveCharacterTextSplitter.from_huggingface_tokenizer(tokenizer, separators=["\n", ".", " "],
                                                                              chunk_size=TOKEN_LIMIT, chunk_overlap=0)
    for line in transcript:
        dialogue = dialogue + line + "\n"
    chunks = text_splitter.split_text(dialogue)
    for chunk in chunks:
        tokens = len(tokenizer.encode(chunk))
        print(f"Chunk {tokens}:\n{chunk}\n\n")
        transcript_documents.append(Document(page_content=chunk))
    return transcript_documents
